/***************************************************************
 *
 * Copyright (c) 2015, Menglong TAN <tanmenglong@gmail.com>
 *
 * This program is free software; you can redistribute it
 * and/or modify it under the terms of the GPL licence
 *
 **************************************************************/

/**
 *
 *
 * @file dot.go
 * @author Menglong TAN <tanmenglong@gmail.com>
 * @date Mon Aug 17 22:01:26 2015
 *
 **/

package dag

import (
	"fmt"
	dot "github.com/awalterschulze/gographviz"
	dotparser "github.com/awalterschulze/gographviz/parser"
	"github.com/crackcell/hpipe/log"
	"strings"
)

//===================================================================
// Public APIs
//===================================================================

type DotSerializer struct{}

func NewDotSerializer() *DotSerializer {
	return &DotSerializer{}
}

func (this *DotSerializer) Serialize(dag *DAG) ([]byte, error) {
	// TODO
	return nil, nil
}

func (this *DotSerializer) Deserialize(data []byte) (*DAG, error) {
	ast, err := dotparser.ParseBytes(data)
	if err != nil {
		return nil, err
	}
	graph := dot.NewGraph()
	dot.Analyse(ast, graph)

	p := NewDAG(graph.Name)

	for src, dests := range graph.Edges.SrcToDsts {
		for dest, _ := range dests {

			if orig, ok := p.Jobs[src]; !ok {
				n := dotNameToDAGJob(graph, src)
				n.Post = append(n.Post, dest)
				p.Jobs[src] = n
			} else {
				orig.Post = append(orig.Post, dest)
			}
			if orig, ok := p.Jobs[dest]; !ok {
				n := dotNameToDAGJob(graph, dest)
				n.Prev = append(n.Prev, src)
				p.Jobs[dest] = n
			} else {
				orig.Prev = append(orig.Prev, src)
			}

			if _, ok := p.InDegrees[src]; !ok {
				p.InDegrees[src] = 0
			}
			if orig, ok := p.InDegrees[dest]; !ok {
				p.InDegrees[dest] = 1
			} else {
				p.InDegrees[dest] = orig + 1
			}

		}
	}

	return p, nil
}

//===================================================================
// Private
//===================================================================

func dotNameToDAGJob(graph *dot.Graph, name string) *Job {
	dotJob, ok := graph.Nodes.Lookup[name]
	if !ok {
		panic(fmt.Errorf("no corresponding node"))
	}

	job, err := dotToDAGJob(dotJob)
	if err != nil {
		panic(err)
	}
	return job
}

func dotToDAGJob(node *dot.Node) (*Job, error) {
	p := NewJob()
	p.Name = node.Name
	p.Attrs = dotToDAGAttrs(node.Attrs)
	if v, ok := p.Attrs["nonstrict"]; ok {
		var err error
		if p.NonStrict, err = parseBoolString(v); err != nil {
			return nil, err
		}
	}
	p.Type = getJobTypeFromAttrs(p.Attrs)
	if p.Type == UnknownJob {
		err := fmt.Errorf("unknown job type: %s for %s", p.Attrs["type"], p.Name)
		log.Error(err)
		return nil, err
	}
	return p, nil
}

func parseBoolString(str string) (bool, error) {
	switch strings.ToLower(str) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool value: %s", str)
	}
}

func dotToDAGAttrs(attrs dot.Attrs) Attrs {
	p := NewAttrs()
	for k, v := range attrs {
		p.Set(k, strings.Trim(v, "\""))
	}
	return p
}

func getJobTypeFromAttrs(attrs Attrs) JobType {
	if val, ok := attrs["type"]; !ok {
		return UnknownJob
	} else {
		return ParseJobType(val)
	}
}
