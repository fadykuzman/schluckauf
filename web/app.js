
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
      const reviewString = group.UpdatedAt ? `last reviewd at: ${group.UpdatedAt}` : "Not yet reviewed"

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

    const imagesDiv = document.querySelector('.images-grid')

    imagesDiv.innerHTML = ''
    files.forEach((file, index) => {
      const fileDiv = createImageDiv(file, index)

      imagesDiv.appendChild(fileDiv)
    });

    imagesDiv.addEventListener('click', async (e) => {
      if (e.target.classList.contains('keep-button')) {
        await handleFileAction(e, files, 'keep')
      } else if (e.target.classList.contains('trash-button')) {
        await handleFileAction(e, files, 'trash')
      }
    })

    const backBtn = document.querySelector('.back-to-groups-button');
    backBtn.onclick = loadGroups;
  } catch (error) {
    showError('Failed to load Group')
    console.error(error)
  }

}

async function handleFileAction(e, files, action) {
  const fileDiv = e.target.closest('.image-item')
  const fileId = parseInt(fileDiv.dataset.fileId)
  const file = files.find(f => f.ID === fileId)
  const duplicateImage = fileDiv.querySelector('.duplicate-image')

  await updateFileAction(file, duplicateImage, action)
}

function createImageDiv(file, index) {
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

async function updateFileAction(file, duplicateImage, action) {
  const previousAction = file.Action

  applyActionState(duplicateImage, action);
  file.Action = action

  try {

    await fetchJSON(`/api/groups/${file.GroupID}/files/${file.ID}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        action: action
      })
    })
  } catch (error) {
    applyActionState(duplicateImage, previousAction)
    file.Action = previousAction
    showError(`Failed to ${action} file`)
  }
}

function formatBytes(bytes) {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}

async function loadGroupStatus() {
  try {
    const stats = await fetchJSON("/api/groups/stats")

    const moveToTrashBtn = document.getElementById('move-to-trash-button')
    moveToTrashBtn.disabled = stats.FilesToTrashCount === 0
    const trashCountSpan = document.getElementById('trash-count')
    trashCountSpan.textContent = stats.FilesToTrashCount

    moveToTrashBtn.onclick = async () => {
      try {
        moveToTrashBtn.disabled = true
        trashCountSpan.textContent = 'Processing...'

        await fetchJSON('/api/files/actions/trash', {
          method: 'POST'
        })

        loadGroupStatus();
      } catch (error) {
        console.error(error)
        showError("Failed to move files to trash")
      }
    }

    document.getElementById('pending-count').textContent = stats.Pending
    document.getElementById('decided-count').textContent = stats.Decided



  } catch (error) {
    showError("Failed to load Groups Statistics")
    console.error(error)
  }
}
loadGroupStatus()
loadGroups()
