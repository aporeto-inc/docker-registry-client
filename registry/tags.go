package registry

import (
	"log"
	"net/http"

	digest "github.com/opencontainers/go-digest"
	yaml "gopkg.in/yaml.v2"
)

type tagsResponse struct {
	Tags []string `json:"tags"`
}

// TagListResponse is the Docker registry API response to /v2/<repo>/tags/list
type TagListResponse struct {
	Child    []string                  `json:"child"`
	Manifest map[string]ManifestRecord `json:"manifest"`
	Name     string                    `json:"name"`
	Tags     []string                  `json:"tags"`
}

// ManifestRecord is the contents of a single image's details
type ManifestRecord struct {
	ImageSizeBytes int      `json:"imageSizeBytes"`
	LayerId        string   `json:"layerId"`
	MediaType      string   `json:"mediaType"`
	Tag            []string `json:"tag"`
	TimeCreatedMs  int      `json:"timeCreatedMs"`
	TimeUploadedMs int      `json:"timeUploadedMs"`
}

// TagsAt gives tag at digest value - if multiple, returns first (shouldn't happen tho)
func (reg *Registry) TagsAt(repo string, dig digest.Digest) (tags []string, err error) {
	url := reg.url("/v2/%s/tags/list", repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := reg.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	tagResp := TagListResponse{}
	err = yaml.NewDecoder(resp.Body).Decode(&tagResp)
	if err != nil {
		log.Fatal("failed to decode response: ", err)
	}
	tagList, err := tagResp.Manifest[dig.String()].Tag, nil
	if err != nil {
		log.Fatal("failed to decode manifest: ", err)
	}
	return tagList, nil
}
