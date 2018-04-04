package buildstructs

// This file was autogenerated by https://github.com/GannettDigital/jstransform

type Video struct {
	AdsEnabled        bool   `json:"adsEnabled"`
	AssetDocumentData string `json:"assetDocumentData"`
	AssetGroup        struct {
		Id       string `json:"id"`
		LogoURL  string `json:"logoURL,omitempty"`
		Name     string `json:"name"`
		SiteCode string `json:"siteCode,omitempty"`
		SiteId   string `json:"siteId"`
		SiteName string `json:"siteName"`
		SstsId   string `json:"sstsId,omitempty"`
		Type     string `json:"type,omitempty"`
		URL      string `json:"URL,omitempty"`
	} `json:"assetGroup"`
	AuthoringBehavior      string `json:"authoringBehavior"`
	AuthoringTypeCode      string `json:"authoringTypeCode"`
	AwsPath                string `json:"awsPath"`
	BackfillDate           string `json:"backfillDate,omitempty"`
	BookReviewPageURL      string `json:"bookReviewPageURL,omitempty"`
	BrightcoveAccountId    string `json:"brightcoveAccountId"`
	BrightcoveId           string `json:"brightcoveId"`
	Byline                 string `json:"byline,omitempty"`
	ContentProtectionState string `json:"contentProtectionState"`
	ContentSourceCode      string `json:"contentSourceCode,omitempty"`
	Contributors           []struct {
		Id   string `json:"id"`
		Name string `json:"name"`
	} `json:"contributors"`
	CreateDate        string `json:"createDate"`
	CreateSystem      string `json:"createSystem"`
	CreateUser        string `json:"createUser"`
	Credit            string `json:"credit,omitempty"`
	EmbargoDate       string `json:"embargoDate"`
	EventDate         string `json:"eventDate"`
	ExcludeFromMobile bool   `json:"excludeFromMobile"`
	ExpirationDate    string `json:"expirationDate,omitempty"`
	Fronts            []struct {
		Id              string `json:"id"`
		Name            string `json:"name"`
		RecommendedDate string `json:"recommendedDate"`
		Type            string `json:"type"`
	} `json:"fronts"`
	GannettTracking    string  `json:"gannettTracking,omitempty"`
	Headline           string  `json:"headline,omitempty"`
	HlsURL             string  `json:"hlsURL,omitempty"`
	Id                 string  `json:"id"`
	InitialPublishDate string  `json:"initialPublishDate,omitempty"`
	IsEvergreen        bool    `json:"isEvergreen"`
	Keywords           string  `json:"keywords,omitempty"`
	Length             float64 `json:"length"`
	Links              struct {
		Assets []struct {
			Id        string `json:"id"`
			Overrides struct {
				Caption string `json:"caption,omitempty"`
			} `json:"overrides,omitempty"`
			Position              float64 `json:"position"`
			RelationshipTypeFlags string  `json:"relationshipTypeFlags"`
		} `json:"assets"`
		PhotoId string `json:"photoId,omitempty"`
	} `json:"links"`
	Mp4URL  string `json:"mp4URL"`
	Origin  string `json:"origin"`
	PageURL struct {
		Long  string `json:"long"`
		Short string `json:"short"`
	} `json:"pageURL"`
	PromoBrief            string `json:"promoBrief,omitempty"`
	PropertyDisplayName   string `json:"propertyDisplayName"`
	PropertyId            string `json:"propertyId"`
	PropertyName          string `json:"propertyName"`
	Publication           string `json:"publication"`
	PublishDate           string `json:"publishDate"`
	PublishSystem         string `json:"publishSystem,omitempty"`
	PublishUser           string `json:"publishUser,omitempty"`
	ReaderCommentsEnabled bool   `json:"readerCommentsEnabled"`
	Renditions            []struct {
		AudioOnly    bool    `json:"audioOnly"`
		Codec        string  `json:"codec"`
		Container    string  `json:"container"`
		DisplayName  string  `json:"displayName,omitempty"`
		Duration     float64 `json:"duration"`
		EncodingRate float64 `json:"encodingRate"`
		Height       float64 `json:"height"`
		Size         float64 `json:"size"`
		Type         string  `json:"type"`
		URL          string  `json:"URL"`
		Width        float64 `json:"width"`
	} `json:"renditions,omitempty"`
	SchemaVersion float64 `json:"schemaVersion"`
	ShortHeadline string  `json:"shortHeadline"`
	Source        string  `json:"source,omitempty"`
	Ssts          struct {
		LeafName                  string `json:"leafName,omitempty"`
		Path                      string `json:"path,omitempty"`
		Section                   string `json:"section,omitempty"`
		Storysubject              string `json:"storysubject,omitempty"`
		Subsection                string `json:"subsection,omitempty"`
		Subtopic                  string `json:"subtopic,omitempty"`
		TaxonomyEntityDisplayName string `json:"taxonomyEntityDisplayName,omitempty"`
		Topic                     string `json:"topic,omitempty"`
	} `json:"ssts,omitempty"`
	StatusName      string   `json:"statusName"`
	StoryHighlights []string `json:"storyHighlights"`
	Tags            []struct {
		DateTagged     string  `json:"dateTagged"`
		Id             string  `json:"id"`
		IsPrimary      bool    `json:"isPrimary"`
		Name           string  `json:"name"`
		ParentId       string  `json:"parentId,omitempty"`
		Path           string  `json:"path"`
		RelevanceScore float64 `json:"relevanceScore"`
		TaggingStatus  string  `json:"taggingStatus"`
		TopicType      string  `json:"topicType"`
		Type           string  `json:"type"`
	} `json:"tags,omitempty"`
	Thumbnail  string `json:"thumbnail,omitempty"`
	Title      string `json:"title"`
	Topic      string `json:"topic"`
	Type       string `json:"type"`
	UpdateDate string `json:"updateDate"`
	UpdateUser string `json:"updateUser"`
	VideoStill string `json:"videoStill"`
}
