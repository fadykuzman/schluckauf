# Workflow

## Initial Use

- [x] User opens the application
- [x] The application loads the main page with no data
- [x] The user selects/types in a path to the directory to be scanned
- [x] The user clicks the "Scan" button to initiate the scanning process
- [ ] If the path is not valid, an error message is displayed (should persist and the user can close it manually)
- [x] If the path is valid, the application scans the directory and displays the results in a grid
  - [Possible Future State] Each group is expanded to show individual files. Not collapsible to offer an overview of all the files and decrease the clicking. The use should have an option to blend out the decided groups
- [x] User can select files to delete or keep using mouse clicks or keyboard shortcuts.
- [ ] Stats shows how many groups are decided and how many are left to decide.
- [ ] A "Move To Trash" Shows how many files are selected to be deleted.
- [ ] User clicks on "Move To Trash" button to delete selected files
- [ ] [Possible Future State] A confirmation dialog appears asking the user to confirm the deletion.
- [ ] [Possible Future State] A "Move To Trash" button for each group
- [ ] If moving to trash is successful, the groups that have all their files decided (keep or trash) and only one file or no files left are removed from the grid.
- [ ] If moving to trash is successful, and a group has only one file left, that is still in pending, then should still show in the grid.
