package buildstructs

// This file was autogenerated by https://github.com/GannettDigital/jstransform

type Base struct {
	AdsEnabled        bool   `json:"adsEnabled"`
	AssetDocumentData string `json:"assetDocumentData"`
	AssetGroup        struct {
		ID         string  `json:"id"`
		LogoURL    string  `json:"logoURL"`
		Name       string  `json:"name"`
		PropertyID float64 `json:"propertyId"`
		SiteCode   string  `json:"siteCode"`
		SiteID     string  `json:"siteId"`
		SiteName   string  `json:"siteName"`
		SstsID     string  `json:"sstsId"`
		Title      string  `json:"title"`
		Type       string  `json:"type"`
		URL        string  `json:"URL"`
	} `json:"assetGroup"`
	AuthoringBehavior      string `json:"authoringBehavior"`
	AuthoringTypeCode      string `json:"authoringTypeCode"`
	AwsPath                string `json:"awsPath"`
	BackfillDate           string `json:"backfillDate"`
	BookReviewPageURL      string `json:"bookReviewPageURL"`
	Byline                 string `json:"byline"`
	ContentProtectionState string `json:"contentProtectionState"`
	ContentSourceCode      string `json:"contentSourceCode"`
	Contributors           []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"contributors"`
	CreateDate        string `json:"createDate"`
	CreateSystem      string `json:"createSystem"`
	CreateUser        string `json:"createUser"`
	EmbargoDate       string `json:"embargoDate"`
	EventDate         string `json:"eventDate"`
	ExcludeFromMobile bool   `json:"excludeFromMobile"`
	ExpirationDate    string `json:"expirationDate"`
	Fronts            []struct {
		ID              string `json:"id"`
		Name            string `json:"name"`
		Type            string `json:"type"`
		RecommendedDate string `json:"recommendedDate"`
	} `json:"fronts"`
	Headline           string `json:"headline"`
	ID                 string `json:"id"`
	InitialPublishDate string `json:"initialPublishDate"`
	IsEvergreen        bool   `json:"isEvergreen"`
	Keywords           string `json:"keywords"`
	Links              struct {
		Assets []struct {
			ID        string `json:"id"`
			Overrides struct {
				Caption string `json:"caption"`
			} `json:"overrides"`
			Position              float64 `json:"position"`
			RelationshipTypeFlags string  `json:"relationshipTypeFlags"`
		} `json:"assets"`
		PhotoID string `json:"photoId"`
	} `json:"links"`
	PageURL struct {
		Long  string `json:"long"`
		Short string `json:"short"`
	} `json:"pageURL"`
	PromoBrief            string `json:"promoBrief"`
	PropertyID            string `json:"propertyId"`
	PropertyName          string `json:"propertyName"`
	PropertyDisplayName   string `json:"propertyDisplayName"`
	Publication           string `json:"publication"`
	PublishDate           string `json:"publishDate"`
	PublishSystem         string `json:"publishSystem"`
	PublishUser           string `json:"publishUser"`
	ReaderCommentsEnabled bool   `json:"readerCommentsEnabled"`
	SchemaVersion         int    `json:"schemaVersion"`
	ShortHeadline         string `json:"shortHeadline"`
	Source                string `json:"source"`
	Ssts                  struct {
		LeafName                  string `json:"leafName"`
		Path                      string `json:"path"`
		Section                   string `json:"section"`
		Storysubject              string `json:"storysubject"`
		Subsection                string `json:"subsection"`
		Subtopic                  string `json:"subtopic"`
		TaxonomyEntityDisplayName string `json:"taxonomyEntityDisplayName"`
		Topic                     string `json:"topic"`
	} `json:"ssts"`
	StatusName      string   `json:"statusName"`
	StoryHighlights []string `json:"storyHighlights"`
	Tags            []struct {
		DateTagged     string  `json:"dateTagged"`
		ID             string  `json:"id"`
		IsPrimary      bool    `json:"isPrimary"`
		Name           string  `json:"name"`
		ParentID       string  `json:"parentId"`
		Path           string  `json:"path"`
		RelevanceScore float64 `json:"relevanceScore"`
		TaggingStatus  string  `json:"taggingStatus"`
		TopicType      string  `json:"topicType"`
		Type           string  `json:"type"`
	} `json:"tags"`
	Title      string `json:"title"`
	Topic      string `json:"topic"`
	Type       string `json:"type"`
	UpdateDate string `json:"updateDate"`
	UpdateUser string `json:"updateUser"`
}
