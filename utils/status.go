/*
** Copyright [2013-2016] [Megam Systems]
**
** Licensed under the Apache License, Version 2.0 (the "License");
** you may not use this file except in compliance with the License.
** You may obtain a copy of the License at
**
** http://www.apache.org/licenses/LICENSE-2.0
**
** Unless required by applicable law or agreed to in writing, software
** distributed under the License is distributed on an "AS IS" BASIS,
** WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
** See the License for the specific language governing permissions and
** limitations under the License.
 */

package utils

// Status represents the status of a unit in vertice
type Status string

func (s Status) String() string {
	return string(s)
}

func (s Status) Event_type() string {
	switch s.String() {
	case LAUNCHING:
		return ONEINSTANCELAUNCHINGTYPE
	case LAUNCHED:
		return ONEINSTANCELAUNCHEDTYPE
	case BOOTSTRAPPING:
		return ONEINSTANCEBOOTSTRAPPINGTYPE
	case BOOTSTRAPPED:
		return ONEINSTANCEBOOTSTRAPPEDTYPE
	case STATEUPPING:
		return ONEINSTANCESTATEUPPINGTYPE
	case STATEUPPED:
		return ONEINSTANCESTATEUPPEDTYPE
	case RUNNING:
		return ONEINSTANCERUNNINGTYPE
	case STARTING:
		return ONEINSTANCESTARTINGTYPE
	case STARTED:
		return ONEINSTANCESTARTEDTYPE
	case STOPPING:
		return ONEINSTANCESTOPPINGTYPE
	case STOPPED:
		return ONEINSTANCESTOPPEDTYPE
	case UPGRADED:
		return ONEINSTANCEUPGRADEDTYPE
	case DESTROYING:
		return ONEINSTANCEDESTROYINGTYPE
	case NUKED:
		return ONEINSTANCEDELETEDTYPE
	case SNAPSHOTTING:
		return ONEINSTANCESNAPSHOTTINGTYPE
	case SNAPSHOTTED:
		return ONEINSTANCESNAPSHOTTEDTYPE
	case COOKBOOKDOWNLOADING:
	    return 	COOKBOOKDOWNLOADINGTYPE
	case COOKBOOKDOWNLOADED:
	    return 	COOKBOOKDOWNLOADEDTYPE
	case COOKBOOKFAILURE:
	    return 	COOKBOOKFAILURETYPE
	case CHEFSOLODOWNLOADING:
			return CHEFSOLODOWNLOADING
	case AUTHKEYSUPDATING:
	    return 	AUTHKEYSUPDATINGTYPE
	case AUTHKEYSUPDATED:
			return 	AUTHKEYSUPDATEDTYPE
	case AUTHKEYSFAILURE:
	    return 	AUTHKEYSFAILURETYPE
	case CHEFCONFIGSETUPPING:
			return ONEINSTANCECHEFCONFIGSETUPPING
	case CHEFCONFIGSETUPPED:
			return ONEINSTANCECHEFCONFIGSETUPPED
	case INSTANCEIPSUPDATING:
	    return 	INSTANCEIPSUPDATINGTYPE
	case INSTANCEIPSUPDATED:
		  return 	INSTANCEIPSUPDATEDTYPE
	case INSTANCEIPSFAILURE:
	    return 	INSTANCEIPSFAILURETYPE
	case CONTAINERNETWORKSUCCESS:
	    return 	CONTAINERNETWORKSUCCESSTYPE
	case CONTAINERNETWORKFAILURE:
	    return 	CONTAINERNETWORKFAILURETYPE
	case DNSNETWORKCREATING:
			return ONEINSTANCEDNSCNAMING
	case DNSNETWORKCREATED:
			return ONEINSTANCEDNSCNAMED
	case DNSNETWORKSKIPPED:
			return ONEINSTANCEDNSNETWORKSKIPPED
	case CLONING:
			return ONEINSTANCEGITCLONING
	case CLONED:
			return ONEINSTANCEGITCLONED

	case ERROR:
		return ONEINSTANCEERRORTYPE
	default:
		return "arrgh"
	}
}

func (s Status) Description(name string) string {
	error_common := "Oops something went wrong on"
	switch s.String() {
	case LAUNCHING:
		return "Your " + name + " machine is initializing.."
	case LAUNCHED:
		return "Machine " + name + " was initialized on cloud.."
	case BOOTSTRAPPING:
		return name + " was bootstrapping.."
	case BOOTSTRAPPED:
		return name + " was bootstrapped.."
	case STATEUPPING:
		return name + " is stateupping.."
	case STATEUPPED:
		return name + " is stateupped.."
	case RUNNING:
		return name + " is running.."
	case CHEFSOLODOWNLOADING:
			return "Chefsolo Downloading .."
	case CHEFSOLODOWNLOADED:
			return "Chefsolo Downloaded .."
	case STARTING:
		return "Starting process initializing on " + name + ".."
	case STARTED:
		return name + " was started.."
	case STOPPING:
		return "Stopping process initializing on " + name + ".."
	case STOPPED:
		return name + " was stopped.."
	case UPGRADED:
		return name + " was upgraded.."
	case DESTROYING:
		return "Destroying process initializing on " + name + ".."
	case NUKED:
		return name + " was removed.."
	case SNAPSHOTTING:
		return "Snapshotting process initializing on " + name + ".."
	case SNAPSHOTTED:
		return name + " was snapcreated.."
	case COOKBOOKDOWNLOADING:
	    return "Chef cookbooks are downloading.."
	case COOKBOOKDOWNLOADED:
	    return "Chef cookbooks are successfully downloaded.."
	case COOKBOOKFAILURE:
			return error_common + "Downloading Cookbooks on " + name + ".."
	case CHEFCONFIGSETUPPING:
			return "Chef config setupping .."
	case CHEFCONFIGSETUPPED:
			return "Chef config setupped .."
	case CLONING:
				return "Cloning your git repository .."
	case CLONED:
				return "Cloned your git repository .."
	case DNSNETWORKCREATING:
			return "DNS CNAME creating " + name + ".."
	case DNSNETWORKCREATED:
			return "DNS CNAME created successfully " + name + ".."
	case DNSNETWORKSKIPPED:
			return "DNS CNAME skipped " + name + ".."
	case AUTHKEYSUPDATING:
	    return "SSH keys are updating on your " + name
	case AUTHKEYSUPDATED:
		  return "SSH keys are updated on your " + name
	case AUTHKEYSFAILURE:
	    return error_common + "Updating Ssh keys on " + name + ".."
	case INSTANCEIPSUPDATING:
	    return "Private and public ips are updating on your " + name
	case INSTANCEIPSUPDATED:
			return "Private and public ips are updated on your " + name
	case INSTANCEIPSFAILURE:
	    return error_common + "Updating private and public ips on " + name + ".."
	case CONTAINERNETWORKSUCCESS:
	    return "Private and public ips are updated on your " + name
	case CONTAINERNETWORKFAILURE:
	    return error_common + "Updating private and public ips on " + name + ".."
	case ERROR:
		return "Oops something went wrong on " + name + ".."
	default:
		return "arrgh"
	}
}
