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

function showGroup(id) {
  alert(`clicked group ${id}`)
}

function formatBytes(bytes) {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
}

loadGroups()
