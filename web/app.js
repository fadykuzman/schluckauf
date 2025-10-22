
async function loadGroups() {
  var groupContainer = document.querySelector('.group-detail-view')
  groupContainer.hidden = true

  try {
    const groups = await fetchJSON('/api/groups')

    const container = document.getElementById('groups-container');
    container.innerHTML = '';


    groups.forEach(group => {
      const item = document.createElement('div');
      item.className = 'group-item';
      const reviewString = group.UpdatedAt ? `last reviewed at: ${group.UpdatedAt}` : "Not yet reviewed"

      item.innerHTML = `
      <img src='/api/image?path=${encodeURIComponent(group.ThumbnailPath)}'>
      <div class="group-info">
        ${group.FileCount} files
        (${formatBytes(group.Size)} each)
        <span class="group-status" data-status="${group.Status.toLowerCase()}">${group.Status}</span>
        <span class="group-updated-at">${reviewString}</span>
      </div>
    `;
      item.onclick = () => showGroup(group.ID);
      container.appendChild(item)
    });
  } catch (error) {
    showError('Failed to load group')
    console.error(error)
  }
}

async function showGroup(id) {
  try {
    const files = await fetchJSON(`/api/groups/${id}`);

    var groupContainer = document.querySelector('.group-detail-view')
    groupContainer.hidden = false

    var groupTitle = document.querySelector('.duplicate-group-title')

    groupTitle.textContent = `Duplicate Group ${id}`

    const imagesGrid = document.querySelector('.images-grid')

    imagesGrid.innerHTML = ''
    files.forEach((file, index) => {
      const fileDiv = createFileDiv(file, index)
      imagesGrid.appendChild(fileDiv)
    });

    const backBtn = document.querySelector('.back-to-groups-button');
    backBtn.onclick = () => {
      groupContainer.hidden = true
    }

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

function createFileDiv(file, index) {
  const fileDiv = document.createElement('div')

  fileDiv.dataset.fileId = file.ID
  fileDiv.dataset.groupId = file.GroupID

  createImageElement(file, index, fileDiv)

  const duplicateImage = fileDiv.querySelector(".duplicate-image")

  applyActionState(duplicateImage, file.Action)

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

function createImageElement(file, index, fileDiv) {
  fileDiv.className = 'image-item';

  fileDiv.innerHTML = `
    <div class="duplicate-image">
      <img src="/api/image?path=${encodeURIComponent(file.Path)}" alt="Image ${index + 1}">
      <div class="metadata">
        <div><strong>Size: </strong> ${formatBytes(file.Filesize)}</div>
      </div>
    </div>
      <div><button class="keep-button">Keep</button><button class="trash-button">Trash</button></div>
    `;

  const pathDiv = document.createElement('div')
  pathDiv.innerHTML = '<strong>Path:</strong>'
  pathDiv.append(file.Path)

  const metaDataDiv = fileDiv.querySelector('.metadata')
  metaDataDiv.prepend(pathDiv)

}

function formatBytes(bytes) {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}

async function loadGroupStatus() {
  try {
    const stats = await fetchJSON("/api/groups/stats")

    updateTrashButtonState(stats.FilesToTrashCount)

    document.getElementById('pending-count').textContent = stats.Pending
    document.getElementById('decided-count').textContent = stats.Decided

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

      if (response.MovedCount > 0) {
        showSuccess(`Successfully moved ${response.MovedCount} of ${response.TotalCount} to trash`)
      }

      if (response.FailedCount > 0) {
        showError(`Failed to move ${response.FailedCount} files to trash`)
        console.log(response.Errors)
      }

      if (response.PartialFailures > 0) {
        showWarning(`Moved to trash but database not updated`)
        console.warn(response.Errors)
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


setupTrashButton()
loadGroupStatus()
loadGroups()
// Attach eventListener to images-grid
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
