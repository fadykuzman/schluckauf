async function loadGroups() {
  const response = await fetch('/api/groups')
  const groups = await response.json()

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
}

async function showGroup(id) {
  const response = await fetch(`/api/groups/${id}`);
  const files = await response.json()

  const container = document.getElementById('groups-container')
  container.innerHTML = '<h2>Duplicate Group ' + id + '</h2>';

  const imagesDiv = document.createElement('div');
  imagesDiv.className = 'images-grid';

  files.forEach((file, index) => {
    console.log(file)
    const fileDiv = createImageDiv(file, index)

    imagesDiv.appendChild(fileDiv)
  });

  container.appendChild(imagesDiv);

  const backBtn = document.createElement('button');
  backBtn.textContent = 'Back to Groups';
  backBtn.onclick = loadGroups;
  container.appendChild(backBtn);

}

function createImageDiv(file, index) {
  const fileDiv = document.createElement('div')
  fileDiv.className = 'image-item';

  fileDiv.innerHTML = `
    <div class="duplicate-image">
      <img src="/api/image?path=${encodeURIComponent(file.Path)}" alt="Image ${index + 1}">
      <div class="metadata">
        <div><strong>Path:</strong> ${file.Path}</div>
        <div><strong>Size: </strong> ${formatBytes(file.Filesize)}</div>
      </div>
    </div>
      <div><button class="keep-button">Keep</button><button class="trash-button">Trash</button></div>
    `;

  const duplicateImage = fileDiv.querySelector(".duplicate-image")


  const keepButton = fileDiv.querySelector(".keep-button");
  keepButton.onclick = async () => {
    duplicateImage.classList.remove("to-delete")
    duplicateImage.classList.add("to-keep")

    const response = await fetch(`/api/groups/${file.GroupID}/files/${file.ID}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        action: "keep"
      })
    })

    console.log(response)

  }

  const trashButton = fileDiv.querySelector(".trash-button")
  trashButton.onclick = async () => {
    duplicateImage.classList.remove("to-keep")
    duplicateImage.classList.add("to-delete")


    const response = await fetch(`/api/groups/${file.GroupID}/files/${file.ID}`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        action: "trash"
      })
    })
  }

  return fileDiv
}

function formatBytes(bytes) {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}

loadGroups()
