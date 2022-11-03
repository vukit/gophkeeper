package tui

import (
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (r *TUI) FileBrowser(form *tview.Form, downloadFolder string) *tview.TreeView {
	rootDir, err := os.UserHomeDir()
	if err != nil {
		rootDir = downloadFolder
	}

	root := tview.NewTreeNode(rootDir).SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().SetRoot(root).SetCurrentNode(root)
	tree.SetBorder(true).SetTitle("[ File browser ]").SetTitleAlign(tview.AlignLeft)

	add := func(target *tview.TreeNode, path string) {
		files, err := os.ReadDir(path)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			node := tview.NewTreeNode(file.Name()).
				SetReference(filepath.Join(path, file.Name()))
			if file.IsDir() {
				node.SetColor(tcell.ColorGreen)
			}

			target.AddChild(node)
		}
	}

	add(root, rootDir)

	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return
		}
		children := node.GetChildren()
		if len(children) == 0 {
			path := reference.(string)
			fileInfo, err := os.Stat(path)
			if err != nil {
				return
			}
			if fileInfo.IsDir() {
				add(node, path)
			} else {
				form.GetFormItemByLabel("File").(*tview.InputField).SetText(path)
			}
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	return tree
}
