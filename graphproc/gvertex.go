package graphproc

func (v *Vertex) GetConfig() any {
	if v.Vstage != nil {
		return v.Vstage.GetConfig()
	}
	return nil
}

func (v *Vertex) GetSelect() *SelectCfg {
	if v.Vstage != nil {
		return v.Vstage.GetSelect()
	}
	return nil
}

func (v *Vertex) SetSelect(selcfg *SelectCfg) {
	if v.Vstage != nil {
		v.Vstage.SetSelect(selcfg)
	}
}

func (v *Vertex) SelectEdge(e *Edge) {
	if v.Vstage != nil {
		v.Vstage.SelectEdge(e)
	}
}

func (v *Vertex) SetConfig(config any) {
	if v.Vstage != nil {
		v.Vstage.BuildState(config)
	}
}
