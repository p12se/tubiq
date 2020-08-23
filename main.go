
package main

import (
	"context"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/jessevdk/go-flags"
	"github.com/rivo/tview"
	"os"
	"path/filepath"
)

var opts struct {
	ProjectID string `short:"p" long:"project" env:"PROJECT_ID" description:"Project ID BigQuery"`
}

var revision = "unknown"

// Show a navigable tree view of the current directory.
func main() {
	fmt.Printf("tubiq %s\n", revision)

	ctx := context.Background()

	if _, err := flags.Parse(&opts); err != nil {
		os.Exit(2)
	}

	bq := newBq(ctx, opts.ProjectID)

	rootDir := ""
	root := tview.NewTreeNode(opts.ProjectID).
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// A helper function which adds the files and directories of the given path
	// to the given target node.
	add := func(target *tview.TreeNode, path string) {
		lists, err := bq.list(path)
		if err != nil {
			panic(err)
		}
		for _, meta := range lists {
			node := tview.NewTreeNode(meta.getName()).
				SetReference(filepath.Join(path, meta.getName())).
				SetSelectable(meta.isDataset())

			if meta.isDataset() {
				node.SetColor(tcell.ColorGreen)
			}
			target.AddChild(node)
		}
	}
	add(root, rootDir)

	// If a directory was selected, open it.
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		reference := node.GetReference()
		if reference == nil {
			return // Selecting the root node does nothing.
		}
		children := node.GetChildren()
		if len(children) == 0 {
			// Load and show fi les in this directory.
			path := reference.(string)
			add(node, path)
		} else {
			// Collapse if visible, expand if collapsed.
			node.SetExpanded(!node.IsExpanded())
		}
	})

	if err := tview.NewApplication().SetRoot(tree, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
