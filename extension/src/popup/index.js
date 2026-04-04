/**
 * Popup da extensao AulaFlash
 * Interface simples de gravacao - injetada direto sem build
 */

const state = { isRecording: false, startTime: null };

document.addEventListener('DOMContentLoaded', async () => {
  const status = await chrome.runtime.sendMessage({ action: 'get-status' });
  state.isRecording = status.isRecording || false;
  state.startTime = status.startTime ? new Date(status.startTime) : null;
  render();
});

function render() {
  const app = document.getElementById('app');
  if (!app) return;

  app.innerHTML = `
    <div style="padding: 16px;">
      <div style="text-align: center;">
        <h1 style="font-size: 18px; font-weight: 700; color: #8B5CF6; margin-bottom: 4px;">
          AulaFlash
        </h1>
        <p style="font-size: 11px; color: #888; margin-bottom: 16px;">
          Grava sua aula/reunião em flashcards
        </p>
      </div>

      <div style="text-align: center; margin: 20px 0;">
        <div style="
          width: 80px; height: 80px; border-radius: 50%;
          background: ${state.isRecording ? '#1a1a2e' : '#2d2d44'};
          border: 3px solid ${state.isRecording ? '#EF4444' : '#8B5CF6'};
          display: flex; align-items: center; justify-content: center;
          margin: 0 auto 12px; cursor: pointer;
          transition: all 0.2s ease;
        " id="recordBtn">
          <div style="
            width: ${state.isRecording ? '24px' : '36px'};
            height: 24px; border-radius: ${state.isRecording ? '4px' : '12px'};
            background: ${state.isRecording ? '#EF4444' : '#8B5CF6'};
            transition: all 0.2s ease;
          "></div>
        </div>

        <p style="font-size: 13px; font-weight: 600;">
          ${state.isRecording ? 'Gravando...' : 'Toque para gravar'}
        </p>

        ${state.isRecording ? `
          <p style="font-size: 11px; color: #888; margin-top: 4px;" id="timer">
            ${getElapsed()}
          </p>
          <div style="
            width: 8px; height: 8px; background: #EF4444;
            border-radius: 50%; display: inline-block;
            margin-right: 4px; animation: pulse 1s infinite;
          "><style>
            @keyframes pulse { 0%, 100% { opacity: 1; } 50% { opacity: 0.3; } }
          </style>
        ` : '<p style="font-size: 10px; color: #666; margin-top: 4px;">Captura discreta da aba</p>'}
      </div>

      <div style="border-top: 1px solid #333; padding-top: 12px;">
        <label style="font-size: 11px; color: #888;">Backend URL</label>
        <input id="backendUrl" type="text" value="http://localhost:8081"
          style="width: 100%; padding: 6px 8px; margin-top: 4px; border: 1px solid #444; border-radius: 6px; background: #2d2d44; color: #e8e8e8; font-size: 11px;"
        />
        <label style="font-size: 11px; color: #888; margin-top: 8px; display: block;">Modo</label>
        <select id="mode"
          style="width: 100%; padding: 6px 8px; margin-top: 4px; border: 1px solid #444; border-radius: 6px; background: #2d2d44; color: #e8e8e8; font-size: 11px;"
        >
          <option value="student">Estudante</option>
          <option value="professional">Profissional</option>
        </select>
        <button id="saveConfig"
          style="width: 100%; padding: 8px; margin-top: 12px; background: #8B5CF6; border: none; border-radius: 6px; color: white; font-size: 11px; cursor: pointer;"
        >Salvar Config</button>
      </div>
    </div>
  `;

  // Event listeners
  document.getElementById('recordBtn').addEventListener('click', toggleRecording);
  document.getElementById('saveConfig').addEventListener('click', saveConfig);

  // Carrega config salva
  chrome.storage.local.get(['backendUrl', 'mode']).then(config => {
    const urlInput = document.getElementById('backendUrl');
    const modeSelect = document.getElementById('mode');
    if (urlInput && config.backendUrl) urlInput.value = config.backendUrl;
    if (modeSelect && config.mode) modeSelect.value = config.mode;
  });

  // Timer
  if (state.isRecording) {
    setInterval(updateTimer, 1000);
  }
}

function getElapsed() {
  if (!state.startTime) return '00:00';
  const diff = Math.floor((Date.now() - state.startTime) / 1000);
  const min = String(Math.floor(diff / 60)).padStart(2, '0');
  const sec = String(diff % 60).padStart(2, '0');
  return `${min}:${sec}`;
}

function updateTimer() {
  const timer = document.getElementById('timer');
  if (timer) timer.textContent = getElapsed();
}

async function toggleRecording() {
  if (state.isRecording) {
    const res = await chrome.runtime.sendMessage({ action: 'stop-recording' });
    if (res.success) {
      state.isRecording = false;
      state.startTime = null;
      render();
    }
  } else {
    const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
    const res = await chrome.runtime.sendMessage({ action: 'start-recording', tabId: tab.id });
    if (res.success) {
      state.isRecording = true;
      state.startTime = new Date();
      render();
    } else {
      alert('Erro: ' + res.error);
    }
  }
}

async function saveConfig() {
  const backendUrl = document.getElementById('backendUrl').value;
  const mode = document.getElementById('mode').value;

  await chrome.storage.local.set({ backendUrl, mode });

  const btn = document.getElementById('saveConfig');
  if (btn) {
    btn.textContent = 'Salvo!';
    btn.style.background = '#10B981';
    setTimeout(() => {
      btn.textContent = 'Salvar Config';
      btn.style.background = '#8B5CF6';
    }, 1500);
  }
}