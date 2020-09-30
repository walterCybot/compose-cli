/*
   Copyright 2020 Docker Compose CLI authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package compose

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/docker/compose-cli/api/client"
	"github.com/docker/compose-cli/api/compose"
	"github.com/docker/compose-cli/formatter"
)

func listCommand() *cobra.Command {
	opts := composeOptions{}
	lsCmd := &cobra.Command{
		Use: "ls",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(cmd.Context(), opts)
		},
	}
	addComposeCommonFlags(lsCmd.Flags(), &opts)
	return lsCmd
}

func addComposeCommonFlags(f *pflag.FlagSet, opts *composeOptions) {
	f.StringVarP(&opts.Name, "project-name", "p", "", "Project name")
	f.StringVar(&opts.Format, "format", "", "Format the output. Values: [pretty | json]. (Default: pretty)")
}

func runList(ctx context.Context, opts composeOptions) error {
	c, err := client.New(ctx)
	if err != nil {
		return err
	}
	stackList, err := c.ComposeService().List(ctx, opts.Name)
	if err != nil {
		return err
	}

	view := viewFromStackList(stackList)
	return formatter.Print(view, opts.Format, os.Stdout, func(w io.Writer) {
		for _, stack := range view {
			_, _ = fmt.Fprintf(w, "%s\t%s\n", stack.Name, stack.Status)
		}
	}, "NAME", "STATUS")
}

type stackView struct {
	Name   string
	Status string
}

func viewFromStackList(stackList []compose.Stack) []stackView {
	retList := make([]stackView, len(stackList))
	for i, s := range stackList {
		retList[i] = stackView{
			Name:   s.Name,
			Status: s.Status,
		}
	}
	return retList
}
