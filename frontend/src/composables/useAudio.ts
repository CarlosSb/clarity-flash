import { ref, readonly } from 'vue'

const isRecording = ref(false)
const mediaRecorder = ref<MediaRecorder | null>(null)
const audioChunks = ref<Blob[]>([])
const recordingTime = ref(0)
let timerInterval: ReturnType<typeof setInterval> | null = null

export function useAudio() {
  async function startRecording(): Promise<void> {
    try {
      const stream = await navigator.mediaDevices.getUserMedia({ audio: true })
      const recorder = new MediaRecorder(stream)
      audioChunks.value = []

      recorder.ondataavailable = (event) => {
        if (event.data.size > 0) {
          audioChunks.value.push(event.data)
        }
      }

      recorder.start()
      mediaRecorder.value = recorder
      isRecording.value = true
      recordingTime.value = 0

      timerInterval = setInterval(() => {
        recordingTime.value++
      }, 1000)
    } catch (error) {
      console.error('Failed to start recording:', error)
      throw error
    }
  }

  function stopRecording(): Promise<Blob> {
    return new Promise((resolve, reject) => {
      if (!mediaRecorder.value || !isRecording.value) {
        reject(new Error('No active recording'))
        return
      }

      const recorder = mediaRecorder.value

      recorder.onstop = () => {
        const blob = new Blob(audioChunks.value, { type: 'audio/webm' })
        // Stop all audio tracks
        recorder.stream.getTracks().forEach((track) => track.stop())

        if (timerInterval) {
          clearInterval(timerInterval)
          timerInterval = null
        }

        isRecording.value = false
        mediaRecorder.value = null

        resolve(blob)
      }

      recorder.stop()
    })
  }

  function cancelRecording(): void {
    if (!mediaRecorder.value || !isRecording.value) return

    const recorder = mediaRecorder.value
    recorder.onstop = () => {
      recorder.stream.getTracks().forEach((track) => track.stop())

      if (timerInterval) {
        clearInterval(timerInterval)
        timerInterval = null
      }

      audioChunks.value = []
      isRecording.value = false
      mediaRecorder.value = null
      recordingTime.value = 0
    }

    try {
      recorder.stop()
    } catch (e) {
      // Already stopped
    }
  }

  function formatTime(seconds: number): string {
    const mins = Math.floor(seconds / 60)
    const secs = seconds % 60
    return `${mins.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`
  }

  return {
    isRecording: readonly(isRecording),
    recordingTime: readonly(recordingTime),
    audioChunks: readonly(audioChunks),
    startRecording,
    stopRecording,
    cancelRecording,
    formatTime,
  }
}
