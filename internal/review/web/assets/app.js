const byId = (id) => document.getElementById(id);

const icons = {
  undo: '<svg viewBox="0 0 24 24" aria-hidden="true"><path d="M9 8 5 12l4 4"/><path d="M5 12h8a6 6 0 1 1 0 12" transform="translate(0 -6)"/></svg>',
};

let documentData = null;
let factState = [];
let questionState = [];
let commentsState = [];
let pendingComment = null;
let currentStep = 0;
let highestStep = 0;
let disabledSnapshot = [];
let commentHoverTimer = null;
const themeStorageKey = 'codex-workflow-theme';

function createElement(tag, className, text) {
  const element = document.createElement(tag);
  if (className) element.className = className;
  if (text !== undefined) element.textContent = text;
  return element;
}

function autoSize(textarea) {
  textarea.style.height = 'auto';
  textarea.style.height = `${Math.max(46, textarea.scrollHeight)}px`;
}

function clearTextSelection() {
  const selection = window.getSelection();
  if (selection) selection.removeAllRanges();
  byId('selection-comment').hidden = true;
}

function preferredTheme() {
  try {
    const stored = localStorage.getItem(themeStorageKey);
    if (stored === 'dark' || stored === 'light') return stored;
  } catch (_error) {
    // 로컬 저장소를 사용할 수 없는 브라우저에서는 OS 설정을 따릅니다.
  }
  return window.matchMedia?.('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
}

function applyTheme(theme, rerender = true) {
  document.documentElement.dataset.theme = theme;
  const toggle = byId('theme-toggle');
  const nextTheme = theme === 'dark' ? 'light' : 'dark';
  toggle.textContent = theme === 'dark' ? '☀' : '☾';
  toggle.setAttribute('aria-label', `${nextTheme === 'dark' ? '다크' : '라이트'} 테마로 전환`);
  toggle.title = toggle.getAttribute('aria-label');
  try {
    localStorage.setItem(themeStorageKey, theme);
  } catch (_error) {
    // 저장 실패는 현재 화면의 테마 적용을 막지 않습니다.
  }
  if (rerender && documentData && currentStep === 2) renderPlanDocument();
}

function showStep(index) {
  currentStep = index;
  highestStep = Math.max(highestStep, index);
  ['facts-screen', 'choices-screen', 'plan-screen'].forEach((id, screenIndex) => {
    byId(id).hidden = screenIndex !== index;
  });
  document.querySelectorAll('.step').forEach((step, stepIndex) => {
    step.disabled = stepIndex > highestStep;
    step.classList.toggle('is-active', stepIndex === index);
    step.classList.toggle('is-complete', stepIndex < index);
    if (stepIndex === index) step.setAttribute('aria-current', 'step');
    else step.removeAttribute('aria-current');
  });
  clearTextSelection();
  if (index === 2) preparePlan();
  window.scrollTo({top: 0, behavior: 'smooth'});
}

function renderFacts() {
  const container = byId('facts');
  container.replaceChildren();
  byId('fact-count').textContent = String(factState.length);
  byId('facts-empty').hidden = factState.length !== 0;

  factState.forEach((fact, index) => {
    const card = createElement('article', 'fact-card');
    const undo = createElement('button', 'icon-button');
    undo.type = 'button';
    undo.title = '원래 사실로 되돌리기';
    undo.setAttribute('aria-label', `${fact.id} 원래 내용으로 되돌리기`);
    undo.innerHTML = icons.undo;

    const main = createElement('div', 'fact-main');
    const meta = createElement('div', 'fact-meta');
    meta.append(createElement('span', 'fact-id', fact.id || `F${index + 1}`));
    meta.append(createElement('span', 'status-badge', fact.status || '확인됨'));
    const textarea = createElement('textarea', 'fact-text');
    textarea.value = fact.content;
    textarea.readOnly = true;
    textarea.rows = 1;
    textarea.setAttribute('aria-label', `${fact.id || `F${index + 1}`} 사실 내용`);
    const evidence = createElement('p', 'evidence', fact.evidence ? `근거 · ${fact.evidence}` : '근거 · 계획 문서');
    main.append(meta, textarea, evidence);

    const edit = createElement('button', 'edit-button', '수정');
    edit.type = 'button';
    edit.setAttribute('aria-label', `${fact.id || `F${index + 1}`} 사실 수정`);

    const sync = () => {
      const edited = fact.content !== fact.original;
      card.classList.toggle('is-edited', edited);
      undo.disabled = !edited;
      autoSize(textarea);
    };
    edit.addEventListener('click', () => {
      const editing = textarea.readOnly;
      textarea.readOnly = !editing;
      edit.textContent = editing ? '완료' : '수정';
      edit.setAttribute(
        'aria-label',
        `${fact.id || `F${index + 1}`} 사실 ${editing ? '편집 완료' : '수정'}`,
      );
      if (editing) {
        textarea.focus();
        textarea.setSelectionRange(textarea.value.length, textarea.value.length);
      }
    });
    textarea.addEventListener('input', () => {
      fact.content = textarea.value;
      sync();
    });
    undo.addEventListener('click', () => {
      fact.content = fact.original;
      textarea.value = fact.original;
      sync();
    });

    card.append(undo, main, edit);
    container.append(card);
    sync();
  });
}

function selectedValues(question) {
  return question.selected.flatMap((value) => {
    if (value !== '__custom__') return [value];
    return question.customValues.map((customValue) => customValue.trim()).filter(Boolean);
  });
}

function optionIsSelected(question, value) {
  return question.selected.includes(value);
}

function defaultModelReason(model) {
  return model === 'gpt-5.6-sol'
    ? '복잡한 구현과 정밀한 검증이 필요해 선택했습니다.'
    : '일반 구현 작업에서 속도와 품질의 균형을 위해 선택했습니다.';
}

function toggleOption(question, value) {
  if (!question.multiple) {
    question.selected = [value];
    return;
  }
  if (question.selected.includes(value)) {
    question.selected = question.selected.filter((selected) => selected !== value);
  } else {
    question.selected = [...question.selected, value];
  }
}

function renderQuestions() {
  const container = byId('questions');
  container.replaceChildren();
  byId('question-count').textContent = String(questionState.length);

  if (questionState.length === 0) {
    const empty = createElement('div', 'empty-state');
    empty.append(
      createElement('strong', '', '지금 선택해야 할 분기가 없습니다.'),
      createElement('span', '', '확인된 사실과 기존 조건만으로 계획을 검토할 수 있습니다.')
    );
    container.append(empty);
    return;
  }

  questionState.forEach((question, questionIndex) => {
    const card = createElement('article', 'question-card');
    const header = createElement('div', 'question-header');
    const kicker = createElement('p', 'question-kicker', question.id || `Q${questionIndex + 1}`);
    kicker.append(createElement('span', 'choice-mode', question.multiple ? '복수 선택' : '하나 선택'));
    header.append(kicker, createElement('h2', '', question.prompt));
    if (question.recommended) {
      const recommendation = createElement('div', 'recommendation');
      recommendation.append(createElement('strong', '', 'AI 권장'), createElement('span', '', question.recommended));
      header.append(recommendation);
    }

    const options = createElement('div', 'options');
    const renderOption = (option, optionIndex, custom = false) => {
      const value = custom ? '__custom__' : option.label;
      const row = createElement('div', 'option');
      row.classList.toggle('is-selected', optionIsSelected(question, value));
      const control = createElement('input');
      control.type = question.multiple ? 'checkbox' : 'radio';
      control.name = `question-${questionIndex}`;
      control.value = value;
      control.checked = optionIsSelected(question, value);
      control.id = `question-${questionIndex}-option-${optionIndex}`;

      const copy = createElement('label', 'option-copy');
      copy.htmlFor = control.id;
      copy.append(
        createElement('strong', '', custom ? '직접 입력' : option.label),
        createElement('span', '', custom ? '원하는 방향을 직접 적습니다.' : (option.description || '이 옵션으로 계획을 진행합니다.'))
      );

      const commentToggle = createElement('button', 'comment-toggle', '의견 추가');
      commentToggle.type = 'button';
      const comment = createElement('textarea', 'option-comment');
      comment.rows = 2;
      comment.placeholder = '이 선택에 반영할 조건이나 우려를 적으세요.';
      comment.value = question.comments[value] || '';
      comment.hidden = comment.value === '';
      comment.addEventListener('input', () => { question.comments[value] = comment.value; });
      commentToggle.addEventListener('click', () => {
        comment.hidden = !comment.hidden;
        if (!comment.hidden) comment.focus();
      });

      const select = () => {
        toggleOption(question, value);
        byId('choice-message').textContent = '';
        renderQuestions();
        if (custom && optionIsSelected(question, '__custom__')) {
          const input = document.querySelector(`[data-custom-question="${questionIndex}"]`);
          if (input) input.focus();
        }
      };
      control.addEventListener('change', select);
      row.addEventListener('click', (event) => {
        if (event.target.closest('button, textarea, input, label')) return;
        select();
      });

      row.append(control, copy, commentToggle);
      if (custom) {
        const customEditor = createElement('div', 'custom-editor');
        const customList = createElement('div', 'custom-input-list');
        const selectCustom = () => {
          if (!optionIsSelected(question, '__custom__')) {
            if (!question.multiple) question.selected = ['__custom__'];
            else question.selected = [...question.selected, '__custom__'];
            control.checked = true;
            row.classList.add('is-selected');
          }
        };

        question.customValues.forEach((customValue, customIndex) => {
          const customRow = createElement('div', 'custom-input-row');
          const customInput = createElement('input', 'custom-input');
          customInput.type = 'text';
          customInput.placeholder = '원하는 옵션을 입력하세요.';
          customInput.value = customValue;
          customInput.dataset.customQuestion = String(questionIndex);
          customInput.dataset.customIndex = String(customIndex);
          customInput.setAttribute('aria-label', `직접 입력 옵션 ${customIndex + 1}`);
          customInput.addEventListener('focus', selectCustom);
          customInput.addEventListener('input', () => {
            question.customValues[customIndex] = customInput.value;
            byId('choice-message').textContent = '';
          });
          customRow.append(customInput);

          if (question.multiple && question.customValues.length > 1) {
            const removeCustom = createElement('button', 'custom-remove', '삭제');
            removeCustom.type = 'button';
            removeCustom.setAttribute('aria-label', `직접 입력 옵션 ${customIndex + 1} 삭제`);
            removeCustom.addEventListener('click', () => {
              question.customValues.splice(customIndex, 1);
              renderQuestions();
            });
            customRow.append(removeCustom);
          }
          customList.append(customRow);
        });
        customEditor.append(customList);

        if (question.multiple) {
          const addCustom = createElement('button', 'custom-add', '+ 직접 입력 추가');
          addCustom.type = 'button';
          addCustom.addEventListener('click', () => {
            selectCustom();
            question.customValues.push('');
            renderQuestions();
            const inputs = document.querySelectorAll(`[data-custom-question="${questionIndex}"]`);
            inputs[inputs.length - 1]?.focus();
          });
          customEditor.append(addCustom);
        }
        row.append(customEditor);
      }
      row.append(comment);
      options.append(row);
    };

    question.options.forEach((option, optionIndex) => renderOption(option, optionIndex));
    renderOption({label: '직접 입력', description: ''}, question.options.length, true);
    card.append(header, options);
    container.append(card);
  });
}

function reviewedFacts() {
  return factState.map((fact) => ({
    id: fact.id,
    status: fact.status,
    content: fact.content.trim(),
    evidence: fact.evidence,
  }));
}

function reviewedSelections() {
  return questionState.map((question) => {
    const options = selectedValues(question);
    const optionComments = question.selected
      .map((key) => question.comments[key] ? `${key === '__custom__' ? question.customValues.map((value) => value.trim()).filter(Boolean).join(', ') : key}: ${question.comments[key]}` : '')
      .filter(Boolean);
    return {
      questionId: question.id,
      question: question.prompt,
      option: options.join(', '),
      options,
      comment: optionComments.join(' / '),
      custom: question.selected.includes('__custom__'),
    };
  });
}

function sectionForNode(node) {
  const container = byId('document');
  const element = node.nodeType === Node.ELEMENT_NODE ? node : node.parentElement;
  const headings = [...container.querySelectorAll('h2, h3')];
  let section = documentData.title;
  for (const heading of headings) {
    if (heading === element || heading.contains(element)) return heading.textContent.trim();
    if (heading.compareDocumentPosition(element) & Node.DOCUMENT_POSITION_FOLLOWING) {
      section = heading.textContent.trim();
    } else {
      break;
    }
  }
  return section;
}

function textRangeForQuote(container, quote) {
  if (!quote) return null;
  const walker = document.createTreeWalker(container, NodeFilter.SHOW_TEXT);
  const nodes = [];
  let fullText = '';
  while (walker.nextNode()) {
    nodes.push({node: walker.currentNode, start: fullText.length});
    fullText += walker.currentNode.data;
  }
  const start = fullText.indexOf(quote);
  if (start < 0) return null;
  const end = start + quote.length;
  const startEntry = nodes.find((entry, index) => start >= entry.start && (index === nodes.length - 1 || start < nodes[index + 1].start));
  const endEntry = nodes.find((entry, index) => end > entry.start && (index === nodes.length - 1 || end <= nodes[index + 1].start));
  if (!startEntry || !endEntry) return null;
  const range = document.createRange();
  range.setStart(startEntry.node, start - startEntry.start);
  range.setEnd(endEntry.node, end - endEntry.start);
  return range;
}

function applyCommentHighlights() {
  const container = byId('document');
  container.querySelectorAll('mark.review-comment').forEach((mark) => mark.replaceWith(...mark.childNodes));
  container.normalize();
  container.querySelectorAll('.has-comment').forEach((element) => {
    element.classList.remove('has-comment');
    element.onmouseenter = null;
    element.onmouseleave = null;
    element.onfocus = null;
    element.onblur = null;
  });

  commentsState.filter((comment) => comment.kind === 'inline').forEach((comment) => {
    const range = textRangeForQuote(container, comment.quote);
    if (!range) return;
    const mark = createElement('mark', 'review-comment');
    try {
      range.surroundContents(mark);
      bindCommentHover(mark, comment);
    } catch (_error) {
      // 여러 블록을 가로지른 선택은 기존 코멘트 데이터만 보존합니다.
    }
  });

  commentsState.filter((comment) => comment.kind === 'block').forEach((comment) => {
    const selector = comment.quote.includes('Mermaid') ? '.mermaid-interactive' : '.reviewable-table';
    const anchor = [...container.querySelectorAll(selector)]
      .find((candidate) => candidate.dataset.section === comment.section);
    if (!anchor) return;
    anchor.classList.add('has-comment');
    bindCommentHover(anchor, comment);
  });
}

function showCommentHover(comment, anchor) {
  window.clearTimeout(commentHoverTimer);
  const hover = byId('comment-hover');
  const rect = anchor.getBoundingClientRect();
  byId('comment-hover-kind').textContent = comment.kind === 'block' ? `블록 · ${comment.section}` : comment.kind === 'global' ? '전역 코멘트' : `본문 · ${comment.section}`;
  byId('comment-hover-text').textContent = comment.comment;
  byId('comment-hover-delete').onclick = () => {
    commentsState = commentsState.filter((candidate) => candidate.id !== comment.id);
    hover.hidden = true;
    renderComments();
  };
  hover.style.left = `${Math.min(window.innerWidth - 300, Math.max(12, rect.left))}px`;
  hover.style.top = `${Math.min(window.innerHeight - 150, rect.bottom + 8)}px`;
  hover.hidden = false;
}

function scheduleCommentHoverClose() {
  window.clearTimeout(commentHoverTimer);
  commentHoverTimer = window.setTimeout(() => { byId('comment-hover').hidden = true; }, 180);
}

function bindCommentHover(anchor, comment) {
  anchor.tabIndex = 0;
  anchor.onmouseenter = () => showCommentHover(comment, anchor);
  anchor.onmouseleave = scheduleCommentHoverClose;
  anchor.onfocus = () => showCommentHover(comment, anchor);
  anchor.onblur = scheduleCommentHoverClose;
}

function makeMermaidInteractive(diagram) {
  const svg = diagram.querySelector('svg');
  if (!svg) return;

  const shell = createElement('div', 'mermaid-interactive');
  shell.dataset.section = sectionForNode(diagram);
  const toolbar = createElement('div', 'mermaid-toolbar');
  toolbar.setAttribute('role', 'toolbar');
  toolbar.setAttribute('aria-label', '다이어그램 보기 도구');
  const zoomLabel = createElement('span', 'mermaid-zoom', '100%');
  zoomLabel.setAttribute('aria-live', 'polite');

  let scale = 1;
  let offsetX = 0;
  let offsetY = 0;
  let dragging = false;

  const applyTransform = () => {
    svg.style.transform = `translate(${offsetX}px, ${offsetY}px) scale(${scale})`;
    zoomLabel.textContent = `${Math.round(scale * 100)}%`;
  };
  const reset = () => {
    scale = 1;
    offsetX = 0;
    offsetY = 0;
    applyTransform();
  };
  const zoom = (delta) => {
    scale = Math.min(3, Math.max(.6, Number((scale + delta).toFixed(2))));
    applyTransform();
  };
  const toolButton = (text, label, action) => {
    const button = createElement('button', 'mermaid-tool', text);
    button.type = 'button';
    button.title = label;
    button.setAttribute('aria-label', label);
    button.addEventListener('click', action);
    return button;
  };

  const zoomOut = toolButton('−', '다이어그램 축소', () => zoom(-.2));
  const zoomIn = toolButton('+', '다이어그램 확대', () => zoom(.2));
  const resetButton = toolButton('초기화', '다이어그램 위치와 배율 초기화', reset);
  const commentButton = toolButton('코멘트', '다이어그램 전체에 코멘트', () => {
    openCommentDialog('block', 'Mermaid 다이어그램 전체', sectionForNode(shell));
  });
  const expandButton = toolButton('크게', '다이어그램 크게 보기', () => {
    const expanded = shell.classList.toggle('is-expanded');
    document.body.classList.toggle('mermaid-modal-open', expanded);
    expandButton.textContent = expanded ? '닫기' : '크게';
    expandButton.setAttribute('aria-pressed', String(expanded));
    expandButton.setAttribute('aria-label', expanded ? '다이어그램 크게 보기 닫기' : '다이어그램 크게 보기');
    diagram.focus();
  });
  expandButton.setAttribute('aria-pressed', 'false');
  toolbar.append(commentButton, zoomOut, zoomLabel, zoomIn, resetButton, expandButton);

  diagram.classList.add('mermaid-viewport');
  diagram.tabIndex = 0;
  diagram.setAttribute('role', 'group');
  diagram.setAttribute('aria-label', '상호작용 다이어그램. 마우스 휠로 확대하거나 축소하고 드래그하여 이동할 수 있습니다.');
  diagram.addEventListener('wheel', (event) => {
    event.preventDefault();
    zoom(event.deltaY < 0 ? .1 : -.1);
  }, {passive: false});
  diagram.addEventListener('pointerdown', (event) => {
    if (event.button !== 0) return;
    dragging = true;
    diagram.classList.add('is-panning');
    diagram.setPointerCapture(event.pointerId);
    event.preventDefault();
  });
  diagram.addEventListener('pointermove', (event) => {
    if (!dragging) return;
    offsetX += event.movementX;
    offsetY += event.movementY;
    applyTransform();
  });
  const stopDragging = (event) => {
    if (!dragging) return;
    dragging = false;
    diagram.classList.remove('is-panning');
    if (diagram.hasPointerCapture(event.pointerId)) diagram.releasePointerCapture(event.pointerId);
  };
  diagram.addEventListener('pointerup', stopDragging);
  diagram.addEventListener('pointercancel', stopDragging);
  diagram.addEventListener('dblclick', reset);
  shell.addEventListener('keydown', (event) => {
    const panStep = 24;
    if (event.key === '+' || event.key === '=') zoom(.2);
    else if (event.key === '-') zoom(-.2);
    else if (event.key === '0') reset();
    else if (event.key === 'ArrowLeft') offsetX -= panStep;
    else if (event.key === 'ArrowRight') offsetX += panStep;
    else if (event.key === 'ArrowUp') offsetY -= panStep;
    else if (event.key === 'ArrowDown') offsetY += panStep;
    else if (event.key === 'Escape' && shell.classList.contains('is-expanded')) expandButton.click();
    else return;
    if (event.key.startsWith('Arrow')) applyTransform();
    event.preventDefault();
  });

  diagram.replaceWith(shell);
  shell.append(toolbar, diagram);
  applyTransform();
}

function makeTableReviewable(table) {
  const section = sectionForNode(table);
  const wrapper = createElement('div', 'reviewable-table');
  wrapper.dataset.section = section;
  const toolbar = createElement('div', 'block-toolbar');
  const button = createElement('button', 'block-comment-button', '표 전체 코멘트');
  button.type = 'button';
  button.addEventListener('click', () => openCommentDialog('block', '표 전체', section));
  toolbar.append(button);
  table.replaceWith(wrapper);
  wrapper.append(toolbar, table);
}

function renderMath(container) {
  if (typeof window.renderMathInElement !== 'function') return;
  window.renderMathInElement(container, {
    delimiters: [
      {left: '$$', right: '$$', display: true},
      {left: '$', right: '$', display: false},
      {left: '\\(', right: '\\)', display: false},
      {left: '\\[', right: '\\]', display: true},
    ],
    throwOnError: false,
    errorColor: '#d45d68',
    strict: 'warn',
    trust: false,
  });
}

function mermaidThemeVariables() {
  if (document.documentElement.dataset.theme === 'light') {
    return {
      fontFamily: 'Segoe UI, sans-serif',
      background: '#f7f9fc',
      primaryColor: '#e9edff',
      primaryTextColor: '#202536',
      primaryBorderColor: '#6574eb',
      lineColor: '#667085',
      secondaryColor: '#f1f3f8',
      secondaryTextColor: '#202536',
      tertiaryColor: '#ffffff',
      tertiaryTextColor: '#202536',
      clusterBkg: '#f1f3f8',
      clusterBorder: '#cbd2e1',
      edgeLabelBackground: '#ffffff',
    };
  }
  return {
    fontFamily: 'Segoe UI, sans-serif',
    background: '#0d1320',
    primaryColor: '#1b2540',
    primaryTextColor: '#dce2ef',
    primaryBorderColor: '#7381ff',
    lineColor: '#8792aa',
    secondaryColor: '#151d2e',
    secondaryTextColor: '#dce2ef',
    tertiaryColor: '#202a3f',
    tertiaryTextColor: '#dce2ef',
    clusterBkg: '#151d2e',
    clusterBorder: '#34405a',
    edgeLabelBackground: '#111726',
  };
}

async function renderPlanDocument() {
  const container = byId('document');
  container.innerHTML = documentData.reviewHtml;
  const headings = [...container.querySelectorAll('h2, h3')];
  const toc = byId('toc');
  toc.replaceChildren();
  headings.forEach((heading, index) => {
    heading.id = `plan-section-${index + 1}`;
    const link = createElement('a', '', heading.textContent);
    link.href = `#${heading.id}`;
    link.dataset.level = heading.tagName === 'H3' ? '3' : '2';
    toc.append(link);
  });
  container.querySelectorAll('a').forEach((link) => {
    link.target = '_blank';
    link.rel = 'noreferrer noopener';
  });
  renderMath(container);
  [...container.querySelectorAll('table')].forEach(makeTableReviewable);

  const mermaidBlocks = [...container.querySelectorAll('pre > code.language-mermaid')];
  const diagrams = mermaidBlocks.map((code) => {
    const diagram = createElement('div', 'mermaid', code.textContent);
    code.parentElement.replaceWith(diagram);
    return diagram;
  });
  if (diagrams.length && window.mermaid) {
    mermaid.initialize({
      startOnLoad: false,
      securityLevel: 'strict',
      theme: 'base',
      themeVariables: mermaidThemeVariables(),
    });
    for (const diagram of diagrams) {
      const source = diagram.textContent;
      try {
        await mermaid.parse(source);
        await mermaid.run({nodes: [diagram]});
        makeMermaidInteractive(diagram);
      } catch (_error) {
        diagram.className = 'mermaid-error';
        diagram.textContent = `Mermaid 구문을 렌더링하지 못했습니다.\n\n${source}`;
      }
    }
  }
  applyCommentHighlights();
}

function renderComments() {
  byId('send-feedback').hidden = commentsState.length === 0;
  const blocked = commentsState.length > 0;
  byId('approve').disabled = blocked;
  byId('approve').title = blocked ? '코멘트가 있어 승인할 수 없습니다.' : '';
  const approveWrap = byId('approve-wrap');
  approveWrap.classList.toggle('is-blocked', blocked);
  approveWrap.tabIndex = blocked ? 0 : -1;
  if (blocked) approveWrap.setAttribute('aria-describedby', 'approval-blocked-tooltip');
  else approveWrap.removeAttribute('aria-describedby');

  const globalComment = commentsState.find((comment) => comment.kind === 'global');
  const globalButton = byId('global-comment');
  globalButton.classList.toggle('has-comment', Boolean(globalComment));
  if (globalComment) bindCommentHover(globalButton, globalComment);
  else {
    globalButton.onmouseenter = null;
    globalButton.onmouseleave = null;
    globalButton.onfocus = null;
    globalButton.onblur = null;
  }
  applyCommentHighlights();
}

function preparePlan() {
  byId('ledger-facts').textContent = String(factState.length);
  byId('ledger-choices').textContent = String(questionState.length);
  renderComments();
  renderPlanDocument();
}

function updateModelDisplay() {
  byId('model-display').textContent = byId('model').value;
  const reason = byId('model-reason').value;
  byId('model-reason-display').textContent = reason;
  byId('model-dialog-reason').textContent = reason;
}

function openModelDialog() {
  const model = byId('model').value;
  const preset = ['gpt-5.6-terra', 'gpt-5.6-sol'].includes(model) ? model : 'custom';
  byId('model-preset').value = preset;
  byId('model-custom').value = preset === 'custom' ? model : '';
  byId('model-custom-wrap').hidden = preset !== 'custom';
  byId('model-error').textContent = '';
  byId('model-dialog').showModal();
  byId('model-preset').focus();
}

function closeModelDialog() {
  byId('model-dialog').close();
}

function saveModelSelection() {
  const preset = byId('model-preset').value;
  const model = (preset === 'custom' ? byId('model-custom').value : preset).trim();
  if (!model.startsWith('gpt-5.6-')) {
    byId('model-error').textContent = 'gpt-5.6 계열 모델을 입력하세요.';
    return;
  }
  byId('model').value = model;
  updateModelDisplay();
  closeModelDialog();
}

function validateChoices() {
  for (const question of questionState) {
    if (selectedValues(question).length === 0) {
      byId('choice-message').textContent = `${question.id || '질문'}의 옵션을 하나 이상 선택하거나 직접 입력하세요.`;
      return false;
    }
  }
  return true;
}

function composePrompt(facts, selections) {
  if (!facts.length && !selections.length) return documentData.prompt;
  const lines = ['[사용자 검토 결과]'];
  if (facts.length) {
    lines.push('', '확정한 사실:');
    facts.forEach((fact) => lines.push(`- ${fact.id}: ${fact.content}`));
  }
  if (selections.length) {
    lines.push('', '선택한 방향:');
    selections.forEach((selection) => {
      const comment = selection.comment ? ` (의견: ${selection.comment})` : '';
      lines.push(`- ${selection.questionId} ${selection.question}: ${selection.options.join(', ')}${comment}`);
    });
  }
  lines.push(
    '',
    '실행 설정:',
    `- 모델: ${byId('model').value.trim()}`,
    `- 모델 선택 이유: ${byId('model-reason').value.trim()}`,
    `- 추론 강도: ${byId('effort').value}`,
  );
  lines.push('', '[기존 실행 지시]', documentData.prompt);
  return lines.join('\n');
}

function openCommentDialog(kind, quote = '', section = '') {
  pendingComment = {kind, quote, section};
  byId('comment-kind').textContent = kind === 'global'
    ? '전역 코멘트'
    : kind === 'block' ? `블록 코멘트 · ${section}` : `본문 코멘트 · ${section}`;
  byId('comment-dialog-title').textContent = kind === 'global'
    ? '계획 전체에 의견 남기기'
    : kind === 'block' ? '블록 전체에 의견 남기기' : '선택한 내용에 의견 남기기';
  byId('comment-quote').textContent = quote;
  byId('comment-text').value = '';
  byId('comment-error').textContent = '';
  byId('comment-dialog').showModal();
  byId('comment-text').focus();
}

function closeCommentDialog() {
  byId('comment-dialog').close();
  pendingComment = null;
  clearTextSelection();
}

function saveComment() {
  const text = byId('comment-text').value.trim();
  if (!text) {
    byId('comment-error').textContent = '코멘트 내용을 입력하세요.';
    return;
  }
  commentsState.push({
    id: `C${Date.now()}-${commentsState.length + 1}`,
    kind: pendingComment.kind,
    section: pendingComment.section,
    quote: pendingComment.quote,
    comment: text,
  });
  closeCommentDialog();
  renderComments();
  applyCommentHighlights();
}

function capturePlanSelection() {
  if (currentStep !== 2 || byId('comment-dialog').open) return;
  const selection = window.getSelection();
  if (!selection || selection.isCollapsed || selection.rangeCount === 0) {
    byId('selection-comment').hidden = true;
    return;
  }
  const range = selection.getRangeAt(0);
  const container = byId('document');
  if (!container.contains(range.commonAncestorContainer)) {
    byId('selection-comment').hidden = true;
    return;
  }
  const quote = selection.toString().trim();
  if (quote.length < 2) {
    byId('selection-comment').hidden = true;
    return;
  }
  pendingComment = {
    kind: 'inline',
    quote: quote.slice(0, 10_000),
    section: sectionForNode(range.startContainer),
  };
  const rect = range.getBoundingClientRect();
  const action = byId('selection-comment');
  action.style.left = `${Math.min(window.innerWidth - 145, Math.max(10, rect.left + rect.width / 2 - 55))}px`;
  action.style.top = `${Math.max(10, rect.top - 42)}px`;
  action.hidden = false;
}

function setSubmitting(submitting) {
  if (submitting) {
    disabledSnapshot = [...document.querySelectorAll('button, input, select, textarea')]
      .map((control) => [control, control.disabled]);
    disabledSnapshot.forEach(([control]) => { control.disabled = true; });
    return;
  }
  disabledSnapshot.forEach(([control, disabled]) => { control.disabled = disabled; });
  disabledSnapshot = [];
}

function showToast(text, success = false) {
  const toast = byId('toast');
  toast.textContent = text;
  toast.classList.toggle('is-success', success);
  toast.hidden = false;
}

function showClosingScreen(action) {
  const approved = action === 'approve';
  byId('app-shell').hidden = true;
  byId('toast').hidden = true;
  byId('closing-screen').hidden = false;
  byId('closing-kicker').textContent = approved ? '검토 완료' : '검토 취소';
  byId('closing-title').textContent = approved ? '승인 결과를 전달했습니다.' : '검토를 취소했습니다.';
  byId('closing-message').textContent = '5초 후 이 페이지가 자동으로 닫힙니다.';
  byId('closing-fallback').hidden = true;
  document.title = approved ? '검토 완료' : '검토 취소';
  history.replaceState(null, '', `${location.pathname}${location.search}#closed`);

  const attemptClose = () => {
    window.close();
    window.setTimeout(() => {
      if (document.hidden) return;
      byId('closing-message').textContent = '검토 세션이 종료되었습니다.';
      byId('closing-fallback').hidden = false;
    }, 350);
  };

  let remaining = 5;
  const timer = window.setInterval(() => {
    remaining -= 1;
    byId('closing-message').textContent = remaining > 0
      ? `${remaining}초 후 이 페이지가 자동으로 닫힙니다.`
      : '페이지를 닫고 있습니다…';
    if (remaining === 0) {
      window.clearInterval(timer);
      attemptClose();
    }
  }, 1000);

  byId('close-now').onclick = () => {
    window.clearInterval(timer);
    attemptClose();
  };
  byId('close-now').focus();
}

async function submit(action) {
  const message = byId('message');
  if (action === 'approve' && commentsState.length > 0) {
    showToast('코멘트가 남아 있어 승인할 수 없습니다. 피드백을 보내거나 코멘트를 삭제하세요.');
    return;
  }
  if (action === 'feedback' && commentsState.length === 0) {
    showToast('전송할 코멘트를 하나 이상 추가하세요.');
    return;
  }
  const model = byId('model').value.trim();
  const modelReason = byId('model-reason').value.trim();
  if (action !== 'cancel' && (!model.startsWith('gpt-5.6-') || !modelReason)) {
    const error = !model.startsWith('gpt-5.6-')
      ? 'gpt-5.6 계열 모델을 선택하거나 입력하세요.'
      : '모델 선택 이유를 입력하세요.';
    message.textContent = error;
    showToast(error);
    byId(!model.startsWith('gpt-5.6-') ? 'model' : 'model-reason').focus();
    return;
  }
  if (action !== 'cancel' && !validateChoices()) {
    showStep(1);
    return;
  }
  const progressText = action === 'approve'
    ? '승인 결과를 반영하고 있습니다…'
    : action === 'feedback' ? '코멘트를 LLM에 전달하고 있습니다…' : '검토를 취소하고 있습니다…';
  message.classList.remove('is-success');
  message.textContent = progressText;
  showToast(progressText);
  setSubmitting(true);
  try {
    const facts = reviewedFacts();
    const selections = reviewedSelections();
    const response = await fetch('./api/submit', {
      method: 'POST',
      headers: {'Content-Type': 'application/json'},
      body: JSON.stringify({
        action,
        model,
        modelReason,
        reasoningEffort: byId('effort').value,
        prompt: action === 'cancel' ? documentData.prompt : composePrompt(facts, selections),
        facts,
        selections,
        comments: commentsState,
      }),
    });
    const result = await response.json();
    if (!response.ok) throw new Error(result.error || '검토 결과를 저장하지 못했습니다.');
    if (action === 'approve' || action === 'cancel') {
      showClosingScreen(action);
      return;
    }
    message.classList.add('is-success');
    message.textContent = '피드백을 보냈습니다. LLM이 수정 범위를 판단한 뒤 필요한 단계에서 다시 시작합니다.';
    showToast(message.textContent, true);
  } catch (error) {
    setSubmitting(false);
    renderComments();
    message.textContent = error.message;
    showToast(error.message);
  }
}

async function load() {
  const response = await fetch('./api/document');
  if (!response.ok) throw new Error('프롬프트 문서를 불러오지 못했습니다.');
  documentData = await response.json();
  factState = (documentData.facts || []).map((fact) => ({...fact, original: fact.content}));
  questionState = (documentData.questions || []).map((question) => {
    const recommended = question.recommended.split(',').map((value) => value.trim()).filter(Boolean);
    const recommendedOptions = question.options.filter((option) => recommended.includes(option.label)).map((option) => option.label);
    return {
      ...question,
      selected: recommendedOptions.length ? recommendedOptions : (question.options[0] ? [question.options[0].label] : []),
      customValues: [''],
      comments: {},
    };
  });

  byId('path').textContent = documentData.path;
  byId('path').title = documentData.path;
  byId('document-title').textContent = documentData.title;
  byId('model').value = documentData.model;
  byId('model-reason').value = documentData.modelReason || defaultModelReason(documentData.model);
  updateModelDisplay();
  if (![...byId('effort').options].some((option) => option.value === documentData.reasoningEffort)) {
    byId('effort').add(new Option(documentData.reasoningEffort, documentData.reasoningEffort));
  }
  byId('effort').value = documentData.reasoningEffort;
  renderFacts();
  renderQuestions();
  renderComments();
  byId('loading').hidden = true;
  const startSteps = {facts: 0, choices: 1, plan: 2};
  const startStep = startSteps[documentData.startAt] ?? 0;
  highestStep = startStep;
  showStep(startStep);
}

document.querySelectorAll('.step').forEach((step) => {
  step.addEventListener('click', () => showStep(Number(step.dataset.step)));
});
applyTheme(preferredTheme(), false);
byId('theme-toggle').addEventListener('click', () => {
  applyTheme(document.documentElement.dataset.theme === 'dark' ? 'light' : 'dark');
});
byId('model-picker').addEventListener('click', openModelDialog);
byId('model-close').addEventListener('click', closeModelDialog);
byId('model-cancel').addEventListener('click', closeModelDialog);
byId('model-save').addEventListener('click', saveModelSelection);
byId('model-preset').addEventListener('change', () => {
  const custom = byId('model-preset').value === 'custom';
  byId('model-custom-wrap').hidden = !custom;
  if (custom) byId('model-custom').focus();
});
byId('comment-hover').addEventListener('mouseenter', () => window.clearTimeout(commentHoverTimer));
byId('comment-hover').addEventListener('mouseleave', scheduleCommentHoverClose);
byId('facts-next').addEventListener('click', () => showStep(1));
byId('choices-back').addEventListener('click', () => showStep(0));
byId('choices-next').addEventListener('click', () => { if (validateChoices()) showStep(2); });
byId('plan-back').addEventListener('click', () => showStep(1));
byId('approve').addEventListener('click', () => submit('approve'));
byId('send-feedback').addEventListener('click', () => submit('feedback'));
byId('top-cancel').addEventListener('click', () => submit('cancel'));
byId('global-comment').addEventListener('click', () => openCommentDialog('global', '', '전체 계획'));
byId('selection-comment').addEventListener('click', () => {
  if (pendingComment) openCommentDialog('inline', pendingComment.quote, pendingComment.section);
});
byId('comment-close').addEventListener('click', closeCommentDialog);
byId('comment-cancel').addEventListener('click', closeCommentDialog);
byId('comment-save').addEventListener('click', saveComment);
byId('comment-dialog').addEventListener('close', () => {
  pendingComment = null;
  clearTextSelection();
});
byId('document').addEventListener('mouseup', () => window.setTimeout(capturePlanSelection, 0));
byId('document').addEventListener('keyup', () => window.setTimeout(capturePlanSelection, 0));
document.addEventListener('mousedown', (event) => {
  if (!event.target.closest('#selection-comment, #document')) byId('selection-comment').hidden = true;
});

load().catch((error) => {
  byId('loading').innerHTML = '';
  byId('loading').append(createElement('strong', '', '검토 화면을 열지 못했습니다.'), createElement('p', '', error.message));
});
