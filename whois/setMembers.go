package whois

import (
	"fmt"
	"strconv"
	"strings"
)

func (w whois) GetSetMembers(set string, recursive bool) ([]string, error) {
	var flagRecursive string
	if recursive {
		flagRecursive = ",1"
	}
	str, err := w.Query("!i" + set + flagRecursive)
	if err != nil {
		return nil, err
	}
	if str == "" {
		return make([]string, 0), nil
	}
	members := strings.Split(str, " ")
	return members, nil
}

func (w whois) GetAsSetMembersRecrusive(set string) ([]uint32, error) {
	members, err := w.GetSetMembers(set, true)
	if err != nil {
		return nil, err
	}

	asns := make([]uint32, 0, len(members))
	for _, member := range members {
		if len(member) < 3 {
			return nil, fmt.Errorf("cannot parse %v into ASN: too short", member)
		}
		asn, err := strconv.Atoi(member[2:])
		if err != nil {
			return nil, fmt.Errorf("cannot parse %v into ASN: %v", member, err)
		}
		asns = append(asns, uint32(asn))
	}
	return asns, nil
}
