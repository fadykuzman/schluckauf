function showError(message) {
  const errorDiv = document.createElement('div')
  errorDiv.className = 'error-toast'
  errorDiv.textContent = message
  document.body.appendChild(errorDiv)

  setTimeout(() => errorDiv.remove(), 5000)
}

function showSuccess(message) {
  const successDiv = document.createElement('div')
  successDiv.className = 'success-toast'
  successDiv.textContent = message
  document.body.appendChild(successDiv)

  setTimeout(() => successDiv.remove(), 5000)
}

function showWarning(message) {
  const warningDiv = document.createElement('div')
  warningDiv.className = 'warning-toast'
  warningDiv.textContent = message
  document.body.appendChild(warningDiv)

  setTimeout(() => warningDiv.remove(), 5000)
}

async function fetchJSON(url, options = {}) {
  const response = await fetch(url, options)

  if (!response.ok) {
    throw new Error(`HTTP ${response.status}: ${response.statusText}`)
  }

  return await response.json()
}

function formatBytes(bytes) {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}
