let selectedImageIndex = null
let currentGroupIndex = 0

function selectImage(index) {
  document.querySelectorAll('.image-item').forEach(item => {
    item.classList.remove('selected')
  })

  selectedImageIndex = null

  const groups = document.querySelectorAll('.duplicate-group')
  if (groups.length === 0 || !groups[currentGroupIndex]) {
    return
  }


  const imageItems = groups[currentGroupIndex].querySelectorAll('.image-item')
  if (index >= 1 && index <= imageItems.length) {
    selectedImageIndex = index
    imageItems[index - 1].classList.add('selected')
  }

  updateShortcutHints()
}

function navigateToGroup(direction) {
  const groups = document.querySelectorAll('.duplicate-group')
  if (groups.length === 0) return

  if (direction === 'down') {
    currentGroupIndex = Math.min(currentGroupIndex + 1, groups.length - 1)
  } else if (direction === 'up') {
    currentGroupIndex = Math.max(currentGroupIndex - 1, 0)
  }

  const targetGroup = groups[currentGroupIndex]
  const firstImage = targetGroup.querySelector('.image-item')

  if (firstImage) {
    const images = targetGroup.querySelectorAll('.image-item')
    selectImage(1)

    targetGroup.scrollIntoView({ behavior: 'smooth', block: 'nearest' })
  }
}

function navigateWithinGroup(direction) {
  const groups = document.querySelectorAll('.duplicate-group')
  if (groups.length === 0) {
    return
  }

  const currentGroup = groups[currentGroupIndex]
  const images = currentGroup.querySelectorAll('.image-item')

  if (images.length === 0) {
    return
  }

  let newIndex = selectedImageIndex || 1

  if (direction === 'right') {
    newIndex = Math.min(newIndex + 1, images.length)
  } else if (direction === 'left') {
    newIndex = Math.max(newIndex - 1, 1)
  }

  selectImage(newIndex)

  images[newIndex - 1].scrollIntoView({ behavior: 'smooth', block: 'nearest' })
}

function updateShortcutHints() {
  const hintsDiv = document.getElementById('shortcut-hints')

  if (selectedImageIndex != null) {
    hintsDiv.innerHTML = `
      <strong>Shortcuts:</strong>
      <kbd>←→</kbd> or <kbd>H/L</kbd> Switch Image |
      <kbd>↑↓</kbd> or <kbd>J/K</kbd> Switch Group |
      <kbd>Enter</kbd> Keep |
      <kbd>Space</kbd> Trash |
      <kbd>Esc</kbd> Deselect 
    `
  } else {
    hintsDiv.innerHTML = `
      <strong>Shortcuts:</strong>
      <kbd>↑↓</kbd> or <kbd>J/K</kbd> Navigate Groups |
      <kbd>1-9</kbd> Select Image |
      <kbd>?</kbd> Help
    `
  }
}

function toggleHelpModal() {
  const modal = document.getElementById('help-modal')
  modal.hidden = !modal.hidden
}


async function loadGroups() {
  try {
    const groups = await fetchJSON('/api/groups')

    const groupsContainer = document.getElementById('groups-container');
    groupsContainer.innerHTML = '';

    if (groups === null || groups.length === 0) {
      const item = document.createElement('div');
      item.textContent = "No duplicates found yet"
      groupsContainer.appendChild(item)

    } else {
      for (const group of groups) {
        const duplicateGroupDiv = document.createElement('div');
        duplicateGroupDiv.className = 'duplicate-group';
        const reviewString = group.updatedAt ? `last reviewed at: ${group.updatedAt}` : "Not yet reviewed"

        duplicateGroupDiv.innerHTML = `
          <div class="group-info">
            ${group.imageCount} files
            <span class="group-status" data-status="${group.status.toLowerCase()}">${group.status}</span>
            <span class="group-updated-at">${reviewString}</span>
          </div>
    `;

        const imagesGrid = await createImagesGrid(group.id)
        groupsContainer.appendChild(duplicateGroupDiv)
        duplicateGroupDiv.appendChild(imagesGrid)

      };
    }
  } catch (error) {
    showError('Failed to load group')
    console.error(error)
  }
  updateShortcutHints()
}

async function createImagesGrid(id) {
  try {
    const files = await fetchJSON(`/api/groups/${id}`);

    var imagesGrid = document.createElement('div')
    imagesGrid.className = 'images-grid'

    imagesGrid.innerHTML = ''
    files.forEach((file, index) => {
      const fileDiv = createImageDiv(file, index)
      imagesGrid.appendChild(fileDiv)
    });
    return imagesGrid
  } catch (error) {
    showError('Failed to load Group')
    console.error(error)
  }
}

function createImageDiv(image, index) {
  const fileDiv = document.createElement('div')

  fileDiv.dataset.fileId = image.id
  fileDiv.dataset.groupId = image.groupId

  createImageElement(image, index, fileDiv)

  const duplicateImage = fileDiv.querySelector(".duplicate-image")

  applyActionState(duplicateImage, image.action)

  return fileDiv
}

function createImageElement(image, index, imageDiv) {
  imageDiv.className = 'image-item';

  imageDiv.innerHTML = `
    <div class="duplicate-image">
      <img src="/api/image?path=${encodeURIComponent(image.path)}" alt="Image ${index + 1}">
      <div class="metadata">
        <div><strong>Size: </strong> ${formatBytes(image.imageSize)}</div>
      </div>
    </div>
    <div class="action-buttons">
      <button class="keep-button">Keep</button>
      <button class="trash-button">Trash</button>
    </div>
    `;

  const pathDiv = document.createElement('div')
  pathDiv.innerHTML = '<strong>Path:</strong>'
  pathDiv.append(image.path)

  const metaDataDiv = imageDiv.querySelector('.metadata')
  metaDataDiv.prepend(pathDiv)
}

function applyActionState(element, action) {
  if (action === "trash") {
    element.classList.remove("to-keep")
    element.classList.add("to-trash")
  } else if (action === "keep") {
    element.classList.remove("to-trash")
    element.classList.add("to-keep")
  }
}

async function updateFileActionById(groupId, fileId, action) {
  const file = document.querySelector(`[data-file-id="${fileId}"]`)
  const duplicateImage = file.querySelector(".duplicate-image")
  let previousAction
  const classList = duplicateImage.classList
  if (classList.contains("to-keep")) {
    previousAction = "keep"
  } else if (classList.contains("to-trash")) {
    previousAction = "trash"
  } else {
    previousAction = "pending"
  }
  applyActionState(duplicateImage, action)
  try {

    await fetchJSON(`/api/groups/${groupId}/files/${fileId}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        action: action
      })
    })
    loadGroupsStatus()
  } catch (error) {
    applyActionState(duplicateImage, previousAction)
    showError(`Failed to ${action} file`)
  }
}

async function loadGroupsStatus() {
  try {
    const stats = await fetchJSON("/api/groups/stats")

    updateTrashButtonState(stats.imagesToTrashCount)

    document.getElementById('pending-count').textContent = stats.pending
    document.getElementById('decided-count').textContent = stats.decided

  } catch (error) {
    showError("Failed to load Groups Statistics")
    console.error(error)
  }
}

function setupTrashButton() {
  const moveToTrashBtn = document.getElementById('move-to-trash-button')
  const trashCountSpan = document.getElementById('trash-count')

  moveToTrashBtn.onclick = async () => {
    try {
      moveToTrashBtn.disabled = true
      trashCountSpan.textContent = 'Processing...'

      const response = await fetchJSON('/api/files/actions/trash', {
        method: 'POST'
      })

      if (response.movedCount > 0) {
        showSuccess(`Successfully moved ${response.movedCount} of ${response.totalCount} to trash`)
      }

      if (response.failedCount > 0) {
        showError(`Failed to move ${response.failedCount} files to trash`)
        console.log(response.errors)
      }

      if (response.partialFailures > 0) {
        showWarning(`Moved to trash but database not updated`)
        console.warn(response.errors)
      }

      loadGroups()
      loadGroupsStatus();
    } catch (error) {
      console.error(error)
      showError("Failed to move files to trash")
    }
  }
}

function updateTrashButtonState(count) {
  const moveToTrashBtn = document.getElementById('move-to-trash-button')
  moveToTrashBtn.disabled = count === 0
  const trashCountSpan = document.getElementById('trash-count')
  trashCountSpan.textContent = count
}

function setupKeyboardShortcuts() {
  document.addEventListener('keydown', (e) => {

    if (document.activeElement.tagName === 'INPUT') {
      return
    }

    if (e.key >= '1' && e.key <= '9') {
      e.preventDefault()
      selectImage(parseInt(e.key))
      return
    }

    const key = e.key.toLowerCase()

    if (e.key === 'ArrowDown' || key === 'j') {
      e.preventDefault()
      navigateToGroup('down')
      return
    }

    if (e.key === 'ArrowUp' || key === 'k') {
      e.preventDefault()
      navigateToGroup('up')
      return
    }

    if (e.key === 'ArrowRight' || key === 'l') {
      e.preventDefault()
      navigateWithinGroup('right')
      return
    }

    if (e.key === 'ArrowLeft' || key === 'h') {
      e.preventDefault()
      navigateWithinGroup('left')
      return
    }

    if (e.key === 'Enter' && selectedImageIndex != null) {
      e.preventDefault()
      const selectedItem = document.querySelector('.image-item.selected')
      if (selectedItem) {
        const fileId = parseInt(selectedItem.dataset.fileId)
        const groupId = parseInt(selectedItem.dataset.groupId)
        updateFileActionById(groupId, fileId, "keep")
      }
    } else if (e.key === ' ' && selectedImageIndex != null) {
      e.preventDefault()
      const selectedItem = document.querySelector('.image-item.selected')
      if (selectedItem) {
        const groupId = parseInt(selectedItem.dataset.groupId)
        const fileId = parseInt(selectedItem.dataset.fileId)
        updateFileActionById(groupId, fileId, "trash")
      }
    } else if (key === 'escape') {
      e.preventDefault()
      selectImage(null)
    }
    if (key === '?' || (e.shiftKey && key === '/' || (e.shiftKey && key === 'ß'))) {
      e.preventDefault()
      toggleHelpModal()
      return
    }
  })

}

function setupFileActionButton() {
  const imagesGrid = document.getElementById('groups-container')
  imagesGrid.addEventListener('click', async (e) => {
    if (e.target.classList.contains('keep-button')) {
      const closest = e.target.closest('.image-item')
      const groupId = parseInt(closest.dataset.groupId)
      const fileId = parseInt(closest.dataset.fileId)
      await updateFileActionById(groupId, fileId, 'keep')
    } else if (e.target.classList.contains('trash-button')) {
      const closest = e.target.closest('.image-item')
      const groupId = parseInt(closest.dataset.groupId)
      const fileId = parseInt(closest.dataset.fileId)
      await updateFileActionById(groupId, fileId, 'trash')
    }
    loadGroupsStatus()
  })
}

function setupScanForm() {
  const form = document.getElementById("scan")
  const input = document.getElementById("scan-directory-input")
  const button = document.getElementById("scan-button")

  form.addEventListener('submit', async (e) => {
    e.preventDefault()

    const directory = input.value.trim()
    if (!directory) {
      showError('Please enter a directory path')
      return
    }

    button.disabled = true
    button.textContent = 'Scanning...'

    try {
      const response = await fetchJSON('/api/scan', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ directory: directory })
      })
      showSuccess(response.message || `Found ${response.groupCount} duplicate groups`)

      await loadGroups()
      await loadGroupsStatus()
    } catch (error) {
      showError('Scan failed ' + error.message)
    } finally {
      button.disabled = false
      button.textContent = 'Scan for Duplicates'
    }
  })
}


function setupHelpModalCloseButton() {
  document.querySelector('.help-modal-close').addEventListener('click', toggleHelpModal)

  document.getElementById('help-modal').addEventListener('click', (e) => {
    if (e.target.id === 'help-modal') {
      toggleHelpModal()
    }
  })
}


loadGroups()
setupTrashButton()
loadGroupsStatus()
setupFileActionButton()
setupKeyboardShortcuts()
setupScanForm()
updateShortcutHints()
setupHelpModalCloseButton()
