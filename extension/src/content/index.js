/**
 * Content Script - AulaFlash
 *
 * Injetado na pagina para detectar contexto (reuniao vs aula)
 * e enviar metadados ao popup.
 */

// Detecta qual plataforma esta ativa
function detectPlatform() {
  const url = window.location.hostname;
  if (url.includes('zoom.us')) return 'zoom';
  if (url.includes('meet.google.com')) return 'google-meet';
  if (url.includes('teams.microsoft.com')) return 'teams';
  if (url.includes('youtube.com')) return 'youtube';
  if (url.includes('classroom.google.com')) return 'classroom';
  if (url.includes('canvas.') || url.includes('moodle')) return 'lms';
  return 'unknown';
}

// Envia deteccao ao popup/background
chrome.runtime.sendMessage({
  action: 'platform-detected',
  platform: detectPlatform(),
  url: window.location.href,
  title: document.title,
});

// Listener para requests do popup
chrome.runtime.onMessage.addListener((message, sender, sendResponse) => {
  if (message.action === 'get-page-info') {
    sendResponse({
      platform: detectPlatform(),
      title: document.title,
      url: window.location.href,
    });
  }
});