# Web Security & Code Quality Fixes Checklist

This document tracks security issues, code quality improvements, and missing features identified in `web/app.js`.

## Critical Security Issues

### XSS Vulnerability (Lines 101-109, 112-113)
- [ ] Replace `innerHTML` usage with safer DOM manipulation methods
- [ ] Use `textContent` instead of `innerHTML` for displaying file paths
- [ ] Sanitize any user-controlled data before inserting into DOM
- [ ] Consider using template elements or createElement throughout

**Current vulnerable code:**
```javascript
pathDiv.innerHTML = '<strong>Path:</strong>'
pathDiv.append(file.Path)  // This is OK, but mixed with innerHTML is confusing
```

**Recommendation:** Use createElement or textContent consistently

---

## Critical Bugs

### Memory Leak: Event Listeners (Line 44)
- [ ] Remove event listeners before clearing container with `innerHTML = ''`
- [ ] Consider using AbortController for cleanup
- [ ] Or use single global event listener on persistent container
- [ ] Document event listener lifecycle

**Issue:** Every `showGroup()` call adds a new listener but never removes old ones

---

### attachEventHandler Function Issues (Lines 63-70)

#### Missing Null Checks
- [ ] Add null check for `fileDiv` from `closest('.image-item')`
- [ ] Add NaN check for `fileId` after `parseInt()`
- [ ] Add undefined check for `file` after `find()`
- [ ] Add null check for `duplicateImage` from `querySelector()`

#### Poor Function Design
- [ ] Rename `attachEventHandler` to `handleFileAction` (more accurate)
- [ ] Add error handling and early returns for invalid states
- [ ] Add logging for debugging purposes

**Improved implementation example:**
```javascript
async function handleFileAction(e, files, action) {
  const fileDiv = e.target.closest('.image-item')
  if (!fileDiv) {
    console.error('Could not find parent image item')
    return
  }

  const fileId = parseInt(fileDiv.dataset.fileId, 10)
  if (isNaN(fileId)) {
    console.error('Invalid file ID')
    return
  }

  const file = files.find(f => f.ID === fileId)
  if (!file) {
    console.error(`File ${fileId} not found`)
    return
  }

  const duplicateImage = fileDiv.querySelector('.duplicate-image')
  if (!duplicateImage) {
    console.error('Could not find image element')
    return
  }

  await updateFileAction(file, duplicateImage, action)
}
```

---

## High Priority Issues

### Missing Keyboard Shortcuts (PRIORITY FEATURE)
- [ ] Implement keyboard event listener (global or scoped)
- [ ] Add number keys (1-9) for image selection
- [ ] Add 'K' key for keep action
- [ ] Add 'D' key for trash action
- [ ] Add 'Esc' key to return to group list
- [ ] Add visual selection highlighting
- [ ] Display keyboard shortcut legend in UI
- [ ] Handle invalid selections gracefully

**Note:** According to CLAUDE.md, keyboard shortcuts are the ESSENTIAL feature for 10x speed improvement

---

### Race Conditions in updateFileAction (Lines 120-142)
- [ ] Disable buttons during API request
- [ ] Implement request queue for rapid clicks
- [ ] Add request cancellation for superseded actions
- [ ] Show loading state on buttons

---

### Missing Null Checks Throughout
- [ ] Line 6: Check `document.getElementById('groups-container')` result
- [ ] Line 30: Check `document.getElementById('groups-container')` result
- [ ] Line 80: Check `querySelector` result before using
- [ ] Line 115: Check `querySelector` result before using

---

## Medium Priority Issues

### Inconsistent Error Messages
- [ ] Standardize capitalization (Lines 21 vs 57)
- [ ] Make error messages more descriptive (include IDs, context)
- [ ] Consider error message constants

**Current inconsistency:**
```javascript
showError('Failed to load group')  // lowercase
showError('Failed to load Group')  // uppercase
```

---

### No Loading States
- [ ] Add loading spinner/indicator during API calls
- [ ] Disable interaction during loading
- [ ] Show progress for long operations
- [ ] Add timeout handling for slow requests

---

### No API Response Validation
- [ ] Validate `groups` is an array in `loadGroups()`
- [ ] Validate `files` is an array in `showGroup()`
- [ ] Check for empty responses
- [ ] Validate required fields exist on objects

**Example:**
```javascript
if (!Array.isArray(files) || files.length === 0) {
  showError('No files found in group')
  return
}
```

---

### Inconsistent DOM Manipulation Patterns
- [ ] Choose one approach: `innerHTML` vs `createElement`
- [ ] Document the chosen pattern
- [ ] Refactor existing code to use chosen pattern consistently

**Current mixing:**
- Template literals with innerHTML (lines 12-16, 101-109)
- createElement + manual construction (lines 10, 33, 111)
- Hybrid approaches (lines 111-116)

---

### Hardcoded Magic Strings
- [ ] Define action constants (`'keep'`, `'trash'`, `'pending'`)
- [ ] Define class name constants
- [ ] Define API endpoint constants

**Recommendation:**
```javascript
const ACTION = {
  KEEP: 'keep',
  TRASH: 'trash',
  PENDING: 'pending'
}

const CLASS = {
  TO_KEEP: 'to-keep',
  TO_TRASH: 'to-trash',
  // etc.
}
```

---

## Low Priority / Nice to Have

### Code Organization
- [ ] Group related functions with comments
- [ ] Consider modules: data loading, DOM creation, events, utilities
- [ ] Add JSDoc comments for complex functions

**Suggested sections:**
```javascript
// === Data Loading ===
// === DOM Creation ===
// === Event Handlers ===
// === Utilities ===
```

---

### Global State & Initialization (Line 150)
- [ ] Wrap immediate execution in DOMContentLoaded handler
- [ ] Add try-catch around initialization
- [ ] Create proper init() function

**Better pattern:**
```javascript
async function init() {
  try {
    await loadGroups()
  } catch (error) {
    showError('Failed to initialize application')
    console.error(error)
  }
}

if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', init)
} else {
  init()
}
```

---

### Inline Event Handlers (Line 17)
- [ ] Replace `onclick` with `addEventListener` for consistency
- [ ] Allows multiple handlers if needed later
- [ ] Easier to remove/manage

---

### Performance: Re-rendering
- [ ] Consider if full DOM rebuild is acceptable for your use case
- [ ] For 100+ images, might want diffing or virtual DOM
- [ ] Or just optimize what you have (probably fine)

---

## Positive Aspects to Preserve

✅ Event delegation pattern (line 44) - efficient and correct
✅ Optimistic UI updates with rollback (lines 121-141) - good UX
✅ Separation of `applyActionState` - clean abstraction
✅ Try-catch error handling in async functions
✅ `encodeURIComponent` for URL parameters - security aware

---

## Testing Checklist

Once fixes are implemented, test:

- [ ] XSS: Try file paths with `<script>alert('xss')</script>` in name
- [ ] Memory: Navigate between groups 50 times, check memory usage
- [ ] Null checks: Remove DOM elements in dev tools, trigger actions
- [ ] Keyboard: All shortcuts work in different contexts
- [ ] Race conditions: Rapidly click keep/trash buttons
- [ ] API errors: Disconnect network, check error handling
- [ ] Loading states: Throttle network in dev tools, check UI feedback

---

## Implementation Priority

**Phase 1 (Critical - Do Before Release):**
1. Fix null checks in `attachEventHandler`
2. Fix XSS vulnerability
3. Implement keyboard shortcuts

**Phase 2 (High Priority):**
4. Fix memory leak
5. Add loading states
6. Fix race conditions

**Phase 3 (Code Quality):**
7. Standardize DOM manipulation
8. Add constants for magic strings
9. Add API response validation

**Phase 4 (Polish):**
10. Improve code organization
11. Add proper initialization
12. Performance optimization if needed
