let selectedImageIndex = null

function selectImage(index) {
  document.querySelectorAll('.image-item').forEach(item => {
    item.classList.remove('selected')
  })

  const imageItems = document.querySelectorAll('.image-item')
  if (index >= 1 && index <= imageItems.length) {
    selectedImageIndex = index
    imageItems[index - 1].classList.add('selected')
  }
}


async function loadGroups() {
  var groupContainer = document.querySelector('.group-detail-view')
  groupContainer.hidden = true

  try {
    const groups = await fetchJSON('/api/groups')

    const container = document.getElementById('groups-container');
    container.innerHTML = '';

    if (groups === null || groups.length === 0) {
      const item = document.createElement('div');
      item.textContent = "No duplicates found yet"
      container.appendChild(item)

    } else {
      for (const group of groups) {
        const item = document.createElement('div');
        item.className = 'duplicate-group';
        const reviewString = group.updatedAt ? `last reviewed at: ${group.updatedAt}` : "Not yet reviewed"

        item.innerHTML = `
          <div class="group-info">
            ${group.imageCount} files
            (${formatBytes(group.size)} each)
            <span class="group-status" data-status="${group.status.toLowerCase()}">${group.status}</span>
            <span class="group-updated-at">${reviewString}</span>
          </div>
    `;

        const imagesGrid = await showGroup(group.id)
        container.appendChild(item)
        container.appendChild(imagesGrid)

      };
    }
  } catch (error) {
    showError('Failed to load group')
    console.error(error)
  }
}

async function showGroup(id) {
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
    loadGroupStatus()
  } catch (error) {
    applyActionState(duplicateImage, previousAction)
    showError(`Failed to ${action} file`)
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

function applyActionState(element, action) {
  if (action === "trash") {
    element.classList.remove("to-keep")
    element.classList.add("to-trash")
  } else if (action === "keep") {
    element.classList.remove("to-trash")
    element.classList.add("to-keep")
  }
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
      <div><button class="keep-button">Keep</button><button class="trash-button">Trash</button></div>
    `;

  const pathDiv = document.createElement('div')
  pathDiv.innerHTML = '<strong>Path:</strong>'
  pathDiv.append(image.path)

  const metaDataDiv = imageDiv.querySelector('.metadata')
  metaDataDiv.prepend(pathDiv)

}


async function loadGroupStatus() {
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

      loadGroupStatus();
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
    const detailView = document.querySelector('.group-detail-view')

    if (detailView.hidden) return

    if (e.key >= '1' && e.key <= '9') {
      e.preventDefault()
      selectImage(parseInt(e.key))
      return
    }

    const key = e.key.toLowerCase()

    if (key === 'k' && selectedImageIndex != null) {
      e.preventDefault()
      const selectedItem = document.querySelector('.image-item.selected')
      if (selectedItem) {
        const fileId = parseInt(selectedItem.dataset.fileId)
        const groupId = parseInt(selectedItem.dataset.groupId)
        updateFileActionById(groupId, fileId, "keep")
      }
    } else if (key === 'd' && selectedImageIndex != null) {
      e.preventDefault()
      const selectedItem = document.querySelector('.image-item.selected')
      if (selectedItem) {
        const groupId = parseInt(selectedItem.dataset.groupId)
        const fileId = parseInt(selectedItem.dataset.fileId)
        updateFileActionById(groupId, fileId, "trash")
      }
    } else if (key === 'escape') {
      e.preventDefault()
      detailView.hidden = true

    }
  })
}

function setupFileActionButton() {
  const imagesGrid = document.querySelector('.images-grid')
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
      await loadGroupStatus()
    } catch (error) {
      showError('Scan failed ' + error.message)
    } finally {
      button.disabled = false
      button.textContent = 'Scan for Duplicates'
    }
  })
}


setupTrashButton()
loadGroupStatus()
loadGroups()
setupFileActionButton()
setupKeyboardShortcuts()
setupScanForm()
