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
    const fileDiv = document.createElement('div')
    fileDiv.className = 'image-item';

    fileDiv.innerHTML = `
      <img src="/api/image?path=${encodeURIComponent(file.Path)}" alt="Image ${index + 1}">
      <div class="metadata">
        <div><strong>Path:</strong> ${file.Path}</div>
        <div><strong>Size: </strong> ${formatBytes(file.Filesize)}</div>
      </div>
    `;


    imagesDiv.appendChild(fileDiv)

  });

  container.appendChild(imagesDiv);

  const backBtn = document.createElement('button');
  backBtn.textContent = 'Back to Groups';
  backBtn.onclick = loadGroups;
  container.appendChild(backBtn);

}

function formatBytes(bytes) {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}

loadGroups()
