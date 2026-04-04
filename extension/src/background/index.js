/**
 * Background Service Worker - AulaFlash Chrome Extension
 *
 * Gerencia a captura de audio via chrome.tabCapture,
 * grava localmente e envia ao backend quando a gravacao para.
 */

let mediaRecorder = null;
let audioChunks = [];
let captureStream = null;
let isRecording = false;
let recordingStartTime = null;

// Inicia gravacao da aba ativa
async function startRecording(tabId) {
  try {
    // Captura o audio da aba
    captureStream = await chrome.tabCapture.capture({
      audio: true,
      video: false,
    });

    if (!captureStream) {
      throw new Error('Nao foi possivel capturar o audio da aba');
    }

    audioChunks = [];
    mediaRecorder = new MediaRecorder(captureStream, {
      mimeType: 'audio/webm;codecs=opus',
    });

    mediaRecorder.ondataavailable = (event) => {
      if (event.data.size > 0) {
        audioChunks.push(event.data);
      }
    };

    mediaRecorder.onstop = async () => {
      const blob = new Blob(audioChunks, { type: 'audio/webm' });
      await uploadAudio(blob);

      // Limpa o stream
      if (captureStream) {
        captureStream.getTracks().forEach(track => track.stop());
        captureStream = null;
      }
    };

    mediaRecorder.start(5000); // Chunk a cada 5s para evitar perda
    isRecording = true;
    recordingStartTime = new Date();

    await chrome.storage.local.set({
      isRecording: true,
      recordingStartTime: recordingStartTime.toISOString(),
    });

    chrome.action.setBadgeText({ text: 'REC' });
    chrome.action.setBadgeBackgroundColor({ color: '#8B5CF6' });
  } catch (error) {
    console.error('Erro ao iniciar gravacao:', error);
    throw error;
  }
}

async function stopRecording() {
  if (mediaRecorder && mediaRecorder.state !== 'inactive') {
    mediaRecorder.stop();
  }

  isRecording = false;
  await chrome.storage.local.set({ isRecording: false });
  chrome.action.setBadgeText({ text: '' });
}

async function uploadAudio(blob) {
  try {
    const { backendUrl, userId, mode } = await chrome.storage.local.get([
      'backendUrl',
      'userId',
      'mode',
    ]);

    const apiUrl = backendUrl || 'http://localhost:8080';
    const uid = userId || generateUserId();

    // Persist user_id
    if (!userId) {
      await chrome.storage.local.set({ userId: uid });
    }

    const formData = new FormData();
    formData.append('audio', blob, 'recording.webm');
    formData.append('user_id', uid);
    formData.append('mode', mode || 'student');

    const response = await fetch(`${apiUrl}/api/sessions/upload`, {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`Upload failed: ${response.status}`);
    }

    const result = await response.json();
    console.log('Upload concluido:', result);

    await chrome.storage.local.set({
      lastSessionId: result.session_id,
      uploadStatus: 'success',
    });
  } catch (error) {
    console.error('Erro no upload:', error);
    await chrome.storage.local.set({ uploadStatus: 'error' });
  }
}

function generateUserId() {
  return 'user_' + Math.random().toString(36).substr(2, 9) + Date.now().toString(36);
}

// Message listener para comunicacao com popup
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === 'start-recording') {
    startRecording(message.tabId)
      .then(() => sendResponse({ success: true }))
      .catch(error => sendResponse({ success: false, error: error.message }));
    return true;
  }

  if (message.action === 'stop-recording') {
    stopRecording()
      .then(() => sendResponse({ success: true }))
      .catch(error => sendResponse({ success: false, error: error.message }));
    return true;
  }

  if (message.action === 'get-status') {
    sendResponse({ isRecording, startTime: recordingStartTime?.toISOString() });
  }
});

chrome.runtime.onInstalled.addListener(async () => {
  await chrome.storage.local.set({ isRecording: false });
});