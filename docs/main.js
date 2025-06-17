// Xử lý sidebar, render bài học, chạy và format code
let current = 0;
let tasksCollapsed = false;

function renderLessonList() {
  const ul = document.getElementById('lessonList');
  ul.innerHTML = '';
  lessons.forEach((l, i) => {
    const li = document.createElement('li');
    li.textContent = l.title;
    if (i === current) li.className = 'active';
    li.onclick = () => {
      current = i;
      renderLesson();
      renderLessonList();
      if (window.innerWidth <= 1100) closeSidebar();
    };
    ul.appendChild(li);
  });
}
function renderLesson() {
  document.getElementById('output').textContent = '';
  const explainDiv = document.getElementById('lessonExplain');
  
  // Render tasks for the current lesson
  renderTasks();
  
  // Nếu có explain thì hiển thị như cũ
  if (lessons[current].explain) {
    explainDiv.style.display = 'block';
    if (lessons[current].format === 'markdown' || !lessons[current].format) {
      explainDiv.innerHTML = marked.parse(lessons[current].explain);
      if (window.Prism) {
        Prism.highlightAllUnder(explainDiv);
      }
    } else if (lessons[current].format === 'html') {
      explainDiv.innerHTML = lessons[current].explain;
    } else {
      explainDiv.innerText = lessons[current].explain;
    }
  } else if (lessons[current].mdPath) {
    // Nếu có mdPath thì fetch file markdown và render
    explainDiv.style.display = 'block';
    fetch(lessons[current].mdPath)
      .then(res => res.text())
      .then(md => {
        explainDiv.innerHTML = marked.parse(md);
        if (window.Prism) {
          Prism.highlightAllUnder(explainDiv);
        }
      });
  } else {
    explainDiv.style.display = 'none';
    explainDiv.innerHTML = '';
  }
  // Fetch code từ file .go
  fetch(lessons[current].codePath)
    .then(res => res.text())
    .then(code => {
      const codeTextarea = document.getElementById('code');
      codeTextarea.value = code;
      codeTextarea.scrollTop = 0;
      codeTextarea.scrollLeft = 0;
      updateLineNumbers();
      setTimeout(() => codeTextarea.focus(), 100);
    });
}
function runGo() {
  document.getElementById('output').textContent = '';
  let ws = new WebSocket("ws://" + window.location.host + "/ws");
  ws.onopen = function() {
    const code = document.getElementById('code').value;
    ws.send(code);
  };
  ws.onmessage = function(evt) {
    if (evt.data === "__DONE__") {
      ws.close();
    } else {
      document.getElementById('output').textContent += evt.data + "\n";
    }
  };
  ws.onerror = function(e) {
    document.getElementById('output').textContent += "[WebSocket error]\n";
  };
}
function formatGo() {
  const code = document.getElementById('code').value;
  fetch('/format', {
    method: 'POST',
    headers: { 'Content-Type': 'text/plain' },
    body: code
  })
  .then(res => {
    if (!res.ok) throw new Error('Format error');
    return res.text();
  })
  .then(formatted => {
    document.getElementById('code').value = formatted;
    updateLineNumbers(); // Cập nhật số dòng sau khi format code
  })
  .catch(() => alert('Không thể format code.'));
}
function openSidebar() {
  document.getElementById('sidebar').classList.add('mobile');
  document.getElementById('sidebarOverlay').classList.add('active');
}
function closeSidebar() {
  document.getElementById('sidebar').classList.remove('mobile');
  document.getElementById('sidebarOverlay').classList.remove('active');
}
// Xử lý chuyển đổi ngôn ngữ
function handleLanguageChange() {
  const languageSelect = document.getElementById('languageSelect');
  if (languageSelect) {
    languageSelect.addEventListener('change', function() {
      const selectedLang = this.value;
      // Thực hiện chuyển đổi ngôn ngữ tại đây
      console.log(`Đã chọn ngôn ngữ: ${selectedLang}`);
      // TODO: Thực hiện thay đổi ngôn ngữ cho toàn bộ ứng dụng
    });
  }
}

// Khởi tạo sau khi DOM đã sẵn sàng
window.addEventListener('DOMContentLoaded', () => {
  renderLessonList();
  renderLesson();
  handleLanguageChange();
  updateLineNumbers();
});

// Hàm cập nhật số dòng
function updateLineNumbers() {
  const codeText = document.getElementById('code').value;
  const lines = codeText.split('\n');
  const lineCount = lines.length;
  const lineNumbersDiv = document.getElementById('lineNumbers');
  
  // Tạo mảng số dòng
  let lineNumbers = [];
  for (let i = 1; i <= lineCount; i++) {
    lineNumbers.push(i);
  }
  
  // Đảm bảo luôn có ít nhất 10 dòng (cho UI đẹp hơn)
  if (lineCount < 10) {
    for (let i = lineCount + 1; i <= 10; i++) {
      lineNumbers.push(i);
    }
  }
  
  // Cập nhật nội dung
  lineNumbersDiv.innerHTML = lineNumbers.join('<br>');
  
  // Đảm bảo cuộn đồng bộ
  lineNumbersDiv.scrollTop = document.getElementById('code').scrollTop;
}

// Hàm đồng bộ cuộn giữa textarea và số dòng
function syncScroll(textarea) {
  const lineNumbers = document.getElementById('lineNumbers');
  
  // Chỉ cần đồng bộ vị trí cuộn theo chiều dọc
  // vì số dòng chỉ cuộn theo chiều dọc
  lineNumbers.scrollTop = textarea.scrollTop;
}

function renderTasks() {
  const tasksContainer = document.getElementById('lessonTasks');
  const tasksList = tasksContainer.querySelector('.tasks-list');
  
  // Ẩn container nếu không có tasks
  if (!lessons[current].tasks || lessons[current].tasks.length === 0) {
    tasksContainer.style.display = 'none';
    return;
  }
  
  tasksContainer.style.display = 'block';
  tasksList.innerHTML = '';
  
  lessons[current].tasks.forEach(task => {
    const taskEl = document.createElement('div');
    taskEl.className = 'task-item';
    
    taskEl.innerHTML = `
      <div class="task-header">
        <div class="task-title">${task.title}</div>
        <span class="task-level ${task.level}">${task.level}</span>
      </div>
      <div class="task-desc">${task.desc}</div>
      ${task.hint ? `<div class="task-hint">${task.hint}</div>` : ''}
    `;
    
    tasksList.appendChild(taskEl);
  });
}

function toggleTasks() {
  const container = document.getElementById('lessonTasks');
  tasksCollapsed = !tasksCollapsed;
  if (tasksCollapsed) {
    container.classList.add('collapsed');
  } else {
    container.classList.remove('collapsed');
  }
}
