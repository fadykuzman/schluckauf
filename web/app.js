async function loadGroups() {
  try {

    const groups = await fetchJSON('/api/groups')

    const container = document.getElementById('groups-container');
    container.innerHTML = '';

    groups.forEach(group => {
      const item = document.createElement('div');
      item.className = 'group-item';
      item.innerHTML = `
      <strong>Group ${group.ID}</strong>:
      ${group.FileCount} files
      (${formatBytes(group.Size)} each)
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

    const container = document.getElementById('groups-container')
    container.innerHTML = '<h2>Duplicate Group ' + id + '</h2>';

    const imagesDiv = document.createElement('div');
    imagesDiv.className = 'images-grid';

    files.forEach((file, index) => {
      const fileDiv = createImageDiv(file, index)

      imagesDiv.appendChild(fileDiv)
    });

    container.appendChild(imagesDiv);

    const backBtn = document.createElement('button');
    backBtn.textContent = 'Back to Groups';
    backBtn.onclick = loadGroups;
    container.appendChild(backBtn);
  } catch (error) {
    showError('Failed to load Group')
    console.error(error)
  }

}

function createImageDiv(file, index) {
  const fileDiv = document.createElement('div')

  createImageElement(file, index, fileDiv)

  const duplicateImage = fileDiv.querySelector(".duplicate-image")

  applyActionState(duplicateImage, file.Action)

  attachButtonHandlers(fileDiv, duplicateImage, file)
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

function attachButtonHandlers(fileDiv, duplicateImage, file) {
  const keepButton = fileDiv.querySelector(".keep-button");
  keepButton.onclick = async () => {
    await updateFileAction(file, duplicateImage, "keep")
  }

  const trashButton = fileDiv.querySelector(".trash-button")
  trashButton.onclick = async () => {
    await updateFileAction(file, duplicateImage, "trash")
  }
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

loadGroups()
