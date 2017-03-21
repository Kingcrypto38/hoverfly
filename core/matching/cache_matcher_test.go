package matching_test

import (
	"testing"

	"github.com/SpectoLabs/hoverfly/core/cache"
	"github.com/SpectoLabs/hoverfly/core/matching"
	"github.com/SpectoLabs/hoverfly/core/models"
	"github.com/SpectoLabs/hoverfly/core/util"
	. "github.com/onsi/gomega"
)

func Test_CacheMatcher_GetCachedResponse_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	_, err := unit.GetCachedResponse(&models.RequestDetails{})
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_GetAllResponses_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	_, err := unit.GetAllResponses()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_SaveRequestTemplateResponsePair_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	err := unit.SaveRequestTemplateResponsePair(models.RequestDetails{}, nil)
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_SaveRequestTemplateResponsePair_WillSaveWithHeaderMatchFalseIfNoHeadesWereOnTheMatchingTemplate(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewInMemoryCache(),
	}

	err := unit.SaveRequestTemplateResponsePair(models.RequestDetails{}, &models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Destination: &models.RequestFieldMatchers{
				ExactMatch: util.StringToPointer("test"),
			},
		},
	})
	Expect(err).To(BeNil())

	cacheValues, err := unit.RequestCache.Get([]byte("d41d8cd98f00b204e9800998ecf8427e"))
	Expect(err).To(BeNil())

	cachedResponse, err := models.NewCachedResponseFromBytes(cacheValues)
	Expect(err).To(BeNil())

	Expect(cachedResponse.HeaderMatch).To(BeFalse())
}

func Test_CacheMatcher_SaveRequestTemplateResponsePair_WillSaveWithHeaderMatchTrueHeadesWereOnTheMatchingTemplate(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewInMemoryCache(),
	}

	err := unit.SaveRequestTemplateResponsePair(models.RequestDetails{}, &models.RequestTemplateResponsePair{
		RequestTemplate: models.RequestTemplate{
			Headers: map[string][]string{
				"test": []string{"headers"},
			},
		},
	})
	Expect(err).To(BeNil())

	cacheValues, err := unit.RequestCache.Get([]byte("d41d8cd98f00b204e9800998ecf8427e"))
	Expect(err).To(BeNil())

	cachedResponse, err := models.NewCachedResponseFromBytes(cacheValues)
	Expect(err).To(BeNil())

	Expect(cachedResponse.HeaderMatch).To(BeTrue())
}

func Test_CacheMatcher_FlushCache_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	err := unit.FlushCache()
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_PreloadCache_WillReturnErrorIfCacheIsNil(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{}

	err := unit.PreloadCache(models.Simulation{})
	Expect(err).ToNot(BeNil())
	Expect(err.Error()).To(Equal("No cache set"))
}

func Test_CacheMatcher_PreloadCache_WillNotCacheLooseTemplates(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewInMemoryCache(),
	}

	err := unit.PreloadCache(models.Simulation{
		Templates: []models.RequestTemplateResponsePair{
			models.RequestTemplateResponsePair{
				RequestTemplate: models.RequestTemplate{
					Body: &models.RequestFieldMatchers{
						RegexMatch: util.StringToPointer("loose"),
					},
				},
				Response: models.ResponseDetails{
					Status: 200,
					Body:   "body",
				},
			},
		},
	})

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.GetAllKeys()).To(HaveLen(0))
}

func Test_CacheMatcher_PreloadCache_WillPreemptivelyCacheFullExactMatchTemplates(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewInMemoryCache(),
	}

	err := unit.PreloadCache(models.Simulation{
		Templates: []models.RequestTemplateResponsePair{
			models.RequestTemplateResponsePair{
				RequestTemplate: models.RequestTemplate{
					Body: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("body"),
					},
					Destination: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("destination"),
					},
					Method: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("method"),
					},
					Path: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("path"),
					},
					Query: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("query"),
					},
					Scheme: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("scheme"),
					},
				},
				Response: models.ResponseDetails{
					Status: 200,
					Body:   "body",
				},
			},
		},
	})

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.GetAllKeys()).To(HaveLen(1))
}

func Test_CacheMatcher_PreloadCache_WillNotPreemptivelyCacheTemplatesWithoutExactMatches(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewInMemoryCache(),
	}

	err := unit.PreloadCache(models.Simulation{
		Templates: []models.RequestTemplateResponsePair{
			models.RequestTemplateResponsePair{
				RequestTemplate: models.RequestTemplate{
					Destination: &models.RequestFieldMatchers{
						RegexMatch: util.StringToPointer("destination"),
					},
				},
				Response: models.ResponseDetails{
					Status: 200,
					Body:   "body",
				},
			},
		},
	})

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.GetAllKeys()).To(HaveLen(0))
}

func Test_CacheMatcher_PreloadCache_WillCheckAllTemplatesInSimulation(t *testing.T) {
	RegisterTestingT(t)
	unit := matching.CacheMatcher{
		RequestCache: cache.NewInMemoryCache(),
	}

	err := unit.PreloadCache(models.Simulation{
		Templates: []models.RequestTemplateResponsePair{
			models.RequestTemplateResponsePair{
				RequestTemplate: models.RequestTemplate{
					Destination: &models.RequestFieldMatchers{
						RegexMatch: util.StringToPointer("destination"),
					},
				},
				Response: models.ResponseDetails{
					Status: 200,
					Body:   "body",
				},
			},
			models.RequestTemplateResponsePair{
				RequestTemplate: models.RequestTemplate{
					Body: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("body"),
					},
					Destination: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("destination"),
					},
					Headers: map[string][]string{
						"Headers": []string{"value"},
					},
					Method: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("method"),
					},
					Path: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("path"),
					},
					Query: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("query"),
					},
					Scheme: &models.RequestFieldMatchers{
						ExactMatch: util.StringToPointer("scheme"),
					},
				},
				Response: models.ResponseDetails{
					Status: 200,
					Body:   "body",
				},
			},
		},
	})

	Expect(err).To(BeNil())
	Expect(unit.RequestCache.GetAllKeys()).To(HaveLen(1))
}
