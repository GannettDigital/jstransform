package rename

// Code generated by github.com/GannettDigital/msgp DO NOT EDIT.

import (
	"time"

	"github.com/GannettDigital/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *ReallyComplex) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Caption":
			z.Caption, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Credit":
			z.Credit, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Crops":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Crops) >= int(zb0002) {
				z.Crops = (z.Crops)[:zb0002]
			} else {
				z.Crops = make([]struct {
					Height       float64 `json:"height"`
					Name         string  `json:"name"`
					Path         string  `json:"path" description:"full path to the cropped image file"`
					RelativePath string  `json:"relativePath" description:"a long"`
					Width        float64 `json:"width"`
				}, zb0002)
			}
			for za0001 := range z.Crops {
				var zb0003 uint32
				zb0003, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				for zb0003 > 0 {
					zb0003--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "Height":
						z.Crops[za0001].Height, err = dc.ReadFloat64()
						if err != nil {
							return
						}
					case "Name":
						z.Crops[za0001].Name, err = dc.ReadString()
						if err != nil {
							return
						}
					case "Path":
						z.Crops[za0001].Path, err = dc.ReadString()
						if err != nil {
							return
						}
					case "RelativePath":
						z.Crops[za0001].RelativePath, err = dc.ReadString()
						if err != nil {
							return
						}
					case "Width":
						z.Crops[za0001].Width, err = dc.ReadFloat64()
						if err != nil {
							return
						}
					default:
						err = dc.Skip()
						if err != nil {
							return
						}
					}
				}
			}
		case "Cutline":
			z.Cutline, err = dc.ReadString()
			if err != nil {
				return
			}
		case "DatePhotoTaken":
			z.DatePhotoTaken, err = dc.ReadTime()
			if err != nil {
				return
			}
		case "Orientation":
			z.Orientation, err = dc.ReadString()
			if err != nil {
				return
			}
		case "OriginalSize":
			var zb0004 uint32
			zb0004, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			for zb0004 > 0 {
				zb0004--
				field, err = dc.ReadMapKeyPtr()
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "Height":
					z.OriginalSize.Height, err = dc.ReadFloat64()
					if err != nil {
						return
					}
				case "Width":
					z.OriginalSize.Width, err = dc.ReadFloat64()
					if err != nil {
						return
					}
				default:
					err = dc.Skip()
					if err != nil {
						return
					}
				}
			}
		case "Type":
			z.Type, err = dc.ReadString()
			if err != nil {
				return
			}
		case "URL":
			var zb0005 uint32
			zb0005, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			for zb0005 > 0 {
				zb0005--
				field, err = dc.ReadMapKeyPtr()
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "Absolute":
					z.URL.Absolute, err = dc.ReadString()
					if err != nil {
						return
					}
				case "Meta":
					var zb0006 uint32
					zb0006, err = dc.ReadMapHeader()
					if err != nil {
						return
					}
					for zb0006 > 0 {
						zb0006--
						field, err = dc.ReadMapKeyPtr()
						if err != nil {
							return
						}
						switch msgp.UnsafeString(field) {
						case "Description":
							z.URL.Meta.Description, err = dc.ReadString()
							if err != nil {
								return
							}
						case "SiteName":
							z.URL.Meta.SiteName, err = dc.ReadString()
							if err != nil {
								return
							}
						default:
							err = dc.Skip()
							if err != nil {
								return
							}
						}
					}
				case "Publish":
					z.URL.Publish, err = dc.ReadString()
					if err != nil {
						return
					}
				default:
					err = dc.Skip()
					if err != nil {
						return
					}
				}
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *ReallyComplex) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 9
	// write "Caption"
	err = en.Append(0x89, 0xa7, 0x43, 0x61, 0x70, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.Caption)
	if err != nil {
		return
	}
	// write "Credit"
	err = en.Append(0xa6, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74)
	if err != nil {
		return
	}
	err = en.WriteString(z.Credit)
	if err != nil {
		return
	}
	// write "Crops"
	err = en.Append(0xa5, 0x43, 0x72, 0x6f, 0x70, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Crops)))
	if err != nil {
		return
	}
	for za0001 := range z.Crops {
		// map header, size 5
		// write "Height"
		err = en.Append(0x85, 0xa6, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74)
		if err != nil {
			return
		}
		err = en.WriteFloat64(z.Crops[za0001].Height)
		if err != nil {
			return
		}
		// write "Name"
		err = en.Append(0xa4, 0x4e, 0x61, 0x6d, 0x65)
		if err != nil {
			return
		}
		err = en.WriteString(z.Crops[za0001].Name)
		if err != nil {
			return
		}
		// write "Path"
		err = en.Append(0xa4, 0x50, 0x61, 0x74, 0x68)
		if err != nil {
			return
		}
		err = en.WriteString(z.Crops[za0001].Path)
		if err != nil {
			return
		}
		// write "RelativePath"
		err = en.Append(0xac, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x50, 0x61, 0x74, 0x68)
		if err != nil {
			return
		}
		err = en.WriteString(z.Crops[za0001].RelativePath)
		if err != nil {
			return
		}
		// write "Width"
		err = en.Append(0xa5, 0x57, 0x69, 0x64, 0x74, 0x68)
		if err != nil {
			return
		}
		err = en.WriteFloat64(z.Crops[za0001].Width)
		if err != nil {
			return
		}
	}
	// write "Cutline"
	err = en.Append(0xa7, 0x43, 0x75, 0x74, 0x6c, 0x69, 0x6e, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Cutline)
	if err != nil {
		return
	}
	// write "DatePhotoTaken"
	err = en.Append(0xae, 0x44, 0x61, 0x74, 0x65, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x54, 0x61, 0x6b, 0x65, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteTime(z.DatePhotoTaken)
	if err != nil {
		return
	}
	// write "Orientation"
	err = en.Append(0xab, 0x4f, 0x72, 0x69, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.Orientation)
	if err != nil {
		return
	}
	// write "OriginalSize"
	// map header, size 2
	// write "Height"
	err = en.Append(0xac, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x53, 0x69, 0x7a, 0x65, 0x82, 0xa6, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.OriginalSize.Height)
	if err != nil {
		return
	}
	// write "Width"
	err = en.Append(0xa5, 0x57, 0x69, 0x64, 0x74, 0x68)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.OriginalSize.Width)
	if err != nil {
		return
	}
	// write "Type"
	err = en.Append(0xa4, 0x54, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Type)
	if err != nil {
		return
	}
	// write "URL"
	// map header, size 3
	// write "Absolute"
	err = en.Append(0xa3, 0x55, 0x52, 0x4c, 0x83, 0xa8, 0x41, 0x62, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.URL.Absolute)
	if err != nil {
		return
	}
	// write "Meta"
	// map header, size 2
	// write "Description"
	err = en.Append(0xa4, 0x4d, 0x65, 0x74, 0x61, 0x82, 0xab, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.URL.Meta.Description)
	if err != nil {
		return
	}
	// write "SiteName"
	err = en.Append(0xa8, 0x53, 0x69, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.URL.Meta.SiteName)
	if err != nil {
		return
	}
	// write "Publish"
	err = en.Append(0xa7, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68)
	if err != nil {
		return
	}
	err = en.WriteString(z.URL.Publish)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *ReallyComplex) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 9
	// string "Caption"
	o = append(o, 0x89, 0xa7, 0x43, 0x61, 0x70, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Caption)
	// string "Credit"
	o = append(o, 0xa6, 0x43, 0x72, 0x65, 0x64, 0x69, 0x74)
	o = msgp.AppendString(o, z.Credit)
	// string "Crops"
	o = append(o, 0xa5, 0x43, 0x72, 0x6f, 0x70, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Crops)))
	for za0001 := range z.Crops {
		// map header, size 5
		// string "Height"
		o = append(o, 0x85, 0xa6, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74)
		o = msgp.AppendFloat64(o, z.Crops[za0001].Height)
		// string "Name"
		o = append(o, 0xa4, 0x4e, 0x61, 0x6d, 0x65)
		o = msgp.AppendString(o, z.Crops[za0001].Name)
		// string "Path"
		o = append(o, 0xa4, 0x50, 0x61, 0x74, 0x68)
		o = msgp.AppendString(o, z.Crops[za0001].Path)
		// string "RelativePath"
		o = append(o, 0xac, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x76, 0x65, 0x50, 0x61, 0x74, 0x68)
		o = msgp.AppendString(o, z.Crops[za0001].RelativePath)
		// string "Width"
		o = append(o, 0xa5, 0x57, 0x69, 0x64, 0x74, 0x68)
		o = msgp.AppendFloat64(o, z.Crops[za0001].Width)
	}
	// string "Cutline"
	o = append(o, 0xa7, 0x43, 0x75, 0x74, 0x6c, 0x69, 0x6e, 0x65)
	o = msgp.AppendString(o, z.Cutline)
	// string "DatePhotoTaken"
	o = append(o, 0xae, 0x44, 0x61, 0x74, 0x65, 0x50, 0x68, 0x6f, 0x74, 0x6f, 0x54, 0x61, 0x6b, 0x65, 0x6e)
	o = msgp.AppendTime(o, z.DatePhotoTaken)
	// string "Orientation"
	o = append(o, 0xab, 0x4f, 0x72, 0x69, 0x65, 0x6e, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.Orientation)
	// string "OriginalSize"
	// map header, size 2
	// string "Height"
	o = append(o, 0xac, 0x4f, 0x72, 0x69, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x53, 0x69, 0x7a, 0x65, 0x82, 0xa6, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74)
	o = msgp.AppendFloat64(o, z.OriginalSize.Height)
	// string "Width"
	o = append(o, 0xa5, 0x57, 0x69, 0x64, 0x74, 0x68)
	o = msgp.AppendFloat64(o, z.OriginalSize.Width)
	// string "Type"
	o = append(o, 0xa4, 0x54, 0x79, 0x70, 0x65)
	o = msgp.AppendString(o, z.Type)
	// string "URL"
	// map header, size 3
	// string "Absolute"
	o = append(o, 0xa3, 0x55, 0x52, 0x4c, 0x83, 0xa8, 0x41, 0x62, 0x73, 0x6f, 0x6c, 0x75, 0x74, 0x65)
	o = msgp.AppendString(o, z.URL.Absolute)
	// string "Meta"
	// map header, size 2
	// string "Description"
	o = append(o, 0xa4, 0x4d, 0x65, 0x74, 0x61, 0x82, 0xab, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.URL.Meta.Description)
	// string "SiteName"
	o = append(o, 0xa8, 0x53, 0x69, 0x74, 0x65, 0x4e, 0x61, 0x6d, 0x65)
	o = msgp.AppendString(o, z.URL.Meta.SiteName)
	// string "Publish"
	o = append(o, 0xa7, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x73, 0x68)
	o = msgp.AppendString(o, z.URL.Publish)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *ReallyComplex) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Caption":
			z.Caption, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Credit":
			z.Credit, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Crops":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Crops) >= int(zb0002) {
				z.Crops = (z.Crops)[:zb0002]
			} else {
				z.Crops = make([]struct {
					Height       float64 `json:"height"`
					Name         string  `json:"name"`
					Path         string  `json:"path" description:"full path to the cropped image file"`
					RelativePath string  `json:"relativePath" description:"a long"`
					Width        float64 `json:"width"`
				}, zb0002)
			}
			for za0001 := range z.Crops {
				var zb0003 uint32
				zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				for zb0003 > 0 {
					zb0003--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "Height":
						z.Crops[za0001].Height, bts, err = msgp.ReadFloat64Bytes(bts)
						if err != nil {
							return
						}
					case "Name":
						z.Crops[za0001].Name, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "Path":
						z.Crops[za0001].Path, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "RelativePath":
						z.Crops[za0001].RelativePath, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "Width":
						z.Crops[za0001].Width, bts, err = msgp.ReadFloat64Bytes(bts)
						if err != nil {
							return
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							return
						}
					}
				}
			}
		case "Cutline":
			z.Cutline, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "DatePhotoTaken":
			z.DatePhotoTaken, bts, err = msgp.ReadTimeBytes(bts)
			if err != nil {
				return
			}
		case "Orientation":
			z.Orientation, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "OriginalSize":
			var zb0004 uint32
			zb0004, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			for zb0004 > 0 {
				zb0004--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "Height":
					z.OriginalSize.Height, bts, err = msgp.ReadFloat64Bytes(bts)
					if err != nil {
						return
					}
				case "Width":
					z.OriginalSize.Width, bts, err = msgp.ReadFloat64Bytes(bts)
					if err != nil {
						return
					}
				default:
					bts, err = msgp.Skip(bts)
					if err != nil {
						return
					}
				}
			}
		case "Type":
			z.Type, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "URL":
			var zb0005 uint32
			zb0005, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			for zb0005 > 0 {
				zb0005--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "Absolute":
					z.URL.Absolute, bts, err = msgp.ReadStringBytes(bts)
					if err != nil {
						return
					}
				case "Meta":
					var zb0006 uint32
					zb0006, bts, err = msgp.ReadMapHeaderBytes(bts)
					if err != nil {
						return
					}
					for zb0006 > 0 {
						zb0006--
						field, bts, err = msgp.ReadMapKeyZC(bts)
						if err != nil {
							return
						}
						switch msgp.UnsafeString(field) {
						case "Description":
							z.URL.Meta.Description, bts, err = msgp.ReadStringBytes(bts)
							if err != nil {
								return
							}
						case "SiteName":
							z.URL.Meta.SiteName, bts, err = msgp.ReadStringBytes(bts)
							if err != nil {
								return
							}
						default:
							bts, err = msgp.Skip(bts)
							if err != nil {
								return
							}
						}
					}
				case "Publish":
					z.URL.Publish, bts, err = msgp.ReadStringBytes(bts)
					if err != nil {
						return
					}
				default:
					bts, err = msgp.Skip(bts)
					if err != nil {
						return
					}
				}
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *ReallyComplex) Msgsize() (s int) {
	s = 1 + 8 + msgp.StringPrefixSize + len(z.Caption) + 7 + msgp.StringPrefixSize + len(z.Credit) + 6 + msgp.ArrayHeaderSize
	for za0001 := range z.Crops {
		s += 1 + 7 + msgp.Float64Size + 5 + msgp.StringPrefixSize + len(z.Crops[za0001].Name) + 5 + msgp.StringPrefixSize + len(z.Crops[za0001].Path) + 13 + msgp.StringPrefixSize + len(z.Crops[za0001].RelativePath) + 6 + msgp.Float64Size
	}
	s += 8 + msgp.StringPrefixSize + len(z.Cutline) + 15 + msgp.TimeSize + 12 + msgp.StringPrefixSize + len(z.Orientation) + 13 + 1 + 7 + msgp.Float64Size + 6 + msgp.Float64Size + 5 + msgp.StringPrefixSize + len(z.Type) + 4 + 1 + 9 + msgp.StringPrefixSize + len(z.URL.Absolute) + 5 + 1 + 12 + msgp.StringPrefixSize + len(z.URL.Meta.Description) + 9 + msgp.StringPrefixSize + len(z.URL.Meta.SiteName) + 8 + msgp.StringPrefixSize + len(z.URL.Publish)
	return
}

// DecodeMsg implements msgp.Decodable
func (z *TotallySimple) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Contributors":
			var zb0002 uint32
			zb0002, err = dc.ReadArrayHeader()
			if err != nil {
				return
			}
			if cap(z.Contributors) >= int(zb0002) {
				z.Contributors = (z.Contributors)[:zb0002]
			} else {
				z.Contributors = make([]struct {
					ContributorId string `json:"contributorId,omitempty"`
					Id            string `json:"id"`
					Name          string `json:"name"`
				}, zb0002)
			}
			for za0001 := range z.Contributors {
				var zb0003 uint32
				zb0003, err = dc.ReadMapHeader()
				if err != nil {
					return
				}
				for zb0003 > 0 {
					zb0003--
					field, err = dc.ReadMapKeyPtr()
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "ContributorId":
						z.Contributors[za0001].ContributorId, err = dc.ReadString()
						if err != nil {
							return
						}
					case "Id":
						z.Contributors[za0001].Id, err = dc.ReadString()
						if err != nil {
							return
						}
					case "Name":
						z.Contributors[za0001].Name, err = dc.ReadString()
						if err != nil {
							return
						}
					default:
						err = dc.Skip()
						if err != nil {
							return
						}
					}
				}
			}
		case "Height":
			z.Height, err = dc.ReadInt64()
			if err != nil {
				return
			}
		case "SomeDateObj":
			var zb0004 uint32
			zb0004, err = dc.ReadMapHeader()
			if err != nil {
				return
			}
			for zb0004 > 0 {
				zb0004--
				field, err = dc.ReadMapKeyPtr()
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "Dates":
					var zb0005 uint32
					zb0005, err = dc.ReadArrayHeader()
					if err != nil {
						return
					}
					if cap(z.SomeDateObj.Dates) >= int(zb0005) {
						z.SomeDateObj.Dates = (z.SomeDateObj.Dates)[:zb0005]
					} else {
						z.SomeDateObj.Dates = make([]time.Time, zb0005)
					}
					for za0002 := range z.SomeDateObj.Dates {
						z.SomeDateObj.Dates[za0002], err = dc.ReadTime()
						if err != nil {
							return
						}
					}
				default:
					err = dc.Skip()
					if err != nil {
						return
					}
				}
			}
		case "Type":
			z.Type, err = dc.ReadString()
			if err != nil {
				return
			}
		case "Visible":
			z.Visible, err = dc.ReadBool()
			if err != nil {
				return
			}
		case "Width":
			z.Width, err = dc.ReadFloat64()
			if err != nil {
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z *TotallySimple) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 6
	// write "Contributors"
	err = en.Append(0x86, 0xac, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x6f, 0x72, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.Contributors)))
	if err != nil {
		return
	}
	for za0001 := range z.Contributors {
		// map header, size 3
		// write "ContributorId"
		err = en.Append(0x83, 0xad, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x6f, 0x72, 0x49, 0x64)
		if err != nil {
			return
		}
		err = en.WriteString(z.Contributors[za0001].ContributorId)
		if err != nil {
			return
		}
		// write "Id"
		err = en.Append(0xa2, 0x49, 0x64)
		if err != nil {
			return
		}
		err = en.WriteString(z.Contributors[za0001].Id)
		if err != nil {
			return
		}
		// write "Name"
		err = en.Append(0xa4, 0x4e, 0x61, 0x6d, 0x65)
		if err != nil {
			return
		}
		err = en.WriteString(z.Contributors[za0001].Name)
		if err != nil {
			return
		}
	}
	// write "Height"
	err = en.Append(0xa6, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74)
	if err != nil {
		return
	}
	err = en.WriteInt64(z.Height)
	if err != nil {
		return
	}
	// write "SomeDateObj"
	// map header, size 1
	// write "Dates"
	err = en.Append(0xab, 0x53, 0x6f, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x65, 0x4f, 0x62, 0x6a, 0x81, 0xa5, 0x44, 0x61, 0x74, 0x65, 0x73)
	if err != nil {
		return
	}
	err = en.WriteArrayHeader(uint32(len(z.SomeDateObj.Dates)))
	if err != nil {
		return
	}
	for za0002 := range z.SomeDateObj.Dates {
		err = en.WriteTime(z.SomeDateObj.Dates[za0002])
		if err != nil {
			return
		}
	}
	// write "Type"
	err = en.Append(0xa4, 0x54, 0x79, 0x70, 0x65)
	if err != nil {
		return
	}
	err = en.WriteString(z.Type)
	if err != nil {
		return
	}
	// write "Visible"
	err = en.Append(0xa7, 0x56, 0x69, 0x73, 0x69, 0x62, 0x6c, 0x65)
	if err != nil {
		return
	}
	err = en.WriteBool(z.Visible)
	if err != nil {
		return
	}
	// write "Width"
	err = en.Append(0xa5, 0x57, 0x69, 0x64, 0x74, 0x68)
	if err != nil {
		return
	}
	err = en.WriteFloat64(z.Width)
	if err != nil {
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z *TotallySimple) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 6
	// string "Contributors"
	o = append(o, 0x86, 0xac, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x6f, 0x72, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.Contributors)))
	for za0001 := range z.Contributors {
		// map header, size 3
		// string "ContributorId"
		o = append(o, 0x83, 0xad, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x6f, 0x72, 0x49, 0x64)
		o = msgp.AppendString(o, z.Contributors[za0001].ContributorId)
		// string "Id"
		o = append(o, 0xa2, 0x49, 0x64)
		o = msgp.AppendString(o, z.Contributors[za0001].Id)
		// string "Name"
		o = append(o, 0xa4, 0x4e, 0x61, 0x6d, 0x65)
		o = msgp.AppendString(o, z.Contributors[za0001].Name)
	}
	// string "Height"
	o = append(o, 0xa6, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74)
	o = msgp.AppendInt64(o, z.Height)
	// string "SomeDateObj"
	// map header, size 1
	// string "Dates"
	o = append(o, 0xab, 0x53, 0x6f, 0x6d, 0x65, 0x44, 0x61, 0x74, 0x65, 0x4f, 0x62, 0x6a, 0x81, 0xa5, 0x44, 0x61, 0x74, 0x65, 0x73)
	o = msgp.AppendArrayHeader(o, uint32(len(z.SomeDateObj.Dates)))
	for za0002 := range z.SomeDateObj.Dates {
		o = msgp.AppendTime(o, z.SomeDateObj.Dates[za0002])
	}
	// string "Type"
	o = append(o, 0xa4, 0x54, 0x79, 0x70, 0x65)
	o = msgp.AppendString(o, z.Type)
	// string "Visible"
	o = append(o, 0xa7, 0x56, 0x69, 0x73, 0x69, 0x62, 0x6c, 0x65)
	o = msgp.AppendBool(o, z.Visible)
	// string "Width"
	o = append(o, 0xa5, 0x57, 0x69, 0x64, 0x74, 0x68)
	o = msgp.AppendFloat64(o, z.Width)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *TotallySimple) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			return
		}
		switch msgp.UnsafeString(field) {
		case "Contributors":
			var zb0002 uint32
			zb0002, bts, err = msgp.ReadArrayHeaderBytes(bts)
			if err != nil {
				return
			}
			if cap(z.Contributors) >= int(zb0002) {
				z.Contributors = (z.Contributors)[:zb0002]
			} else {
				z.Contributors = make([]struct {
					ContributorId string `json:"contributorId,omitempty"`
					Id            string `json:"id"`
					Name          string `json:"name"`
				}, zb0002)
			}
			for za0001 := range z.Contributors {
				var zb0003 uint32
				zb0003, bts, err = msgp.ReadMapHeaderBytes(bts)
				if err != nil {
					return
				}
				for zb0003 > 0 {
					zb0003--
					field, bts, err = msgp.ReadMapKeyZC(bts)
					if err != nil {
						return
					}
					switch msgp.UnsafeString(field) {
					case "ContributorId":
						z.Contributors[za0001].ContributorId, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "Id":
						z.Contributors[za0001].Id, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					case "Name":
						z.Contributors[za0001].Name, bts, err = msgp.ReadStringBytes(bts)
						if err != nil {
							return
						}
					default:
						bts, err = msgp.Skip(bts)
						if err != nil {
							return
						}
					}
				}
			}
		case "Height":
			z.Height, bts, err = msgp.ReadInt64Bytes(bts)
			if err != nil {
				return
			}
		case "SomeDateObj":
			var zb0004 uint32
			zb0004, bts, err = msgp.ReadMapHeaderBytes(bts)
			if err != nil {
				return
			}
			for zb0004 > 0 {
				zb0004--
				field, bts, err = msgp.ReadMapKeyZC(bts)
				if err != nil {
					return
				}
				switch msgp.UnsafeString(field) {
				case "Dates":
					var zb0005 uint32
					zb0005, bts, err = msgp.ReadArrayHeaderBytes(bts)
					if err != nil {
						return
					}
					if cap(z.SomeDateObj.Dates) >= int(zb0005) {
						z.SomeDateObj.Dates = (z.SomeDateObj.Dates)[:zb0005]
					} else {
						z.SomeDateObj.Dates = make([]time.Time, zb0005)
					}
					for za0002 := range z.SomeDateObj.Dates {
						z.SomeDateObj.Dates[za0002], bts, err = msgp.ReadTimeBytes(bts)
						if err != nil {
							return
						}
					}
				default:
					bts, err = msgp.Skip(bts)
					if err != nil {
						return
					}
				}
			}
		case "Type":
			z.Type, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				return
			}
		case "Visible":
			z.Visible, bts, err = msgp.ReadBoolBytes(bts)
			if err != nil {
				return
			}
		case "Width":
			z.Width, bts, err = msgp.ReadFloat64Bytes(bts)
			if err != nil {
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z *TotallySimple) Msgsize() (s int) {
	s = 1 + 13 + msgp.ArrayHeaderSize
	for za0001 := range z.Contributors {
		s += 1 + 14 + msgp.StringPrefixSize + len(z.Contributors[za0001].ContributorId) + 3 + msgp.StringPrefixSize + len(z.Contributors[za0001].Id) + 5 + msgp.StringPrefixSize + len(z.Contributors[za0001].Name)
	}
	s += 7 + msgp.Int64Size + 12 + 1 + 6 + msgp.ArrayHeaderSize + (len(z.SomeDateObj.Dates) * (msgp.TimeSize)) + 5 + msgp.StringPrefixSize + len(z.Type) + 8 + msgp.BoolSize + 6 + msgp.Float64Size
	return
}
