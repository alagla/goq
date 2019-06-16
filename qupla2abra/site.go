package qupla2abra

func (site *SiteIR) Size() int {
	if site.siteType == SITE_MERGE {
		return site.inputs[0].Size()
	}
	ret := 0
	for _, s := range site.inputs {
		ret += s.Size()
	}
	return ret
}
