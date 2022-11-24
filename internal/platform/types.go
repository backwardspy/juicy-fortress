package platform

// archive types "enum" :(
const (
	ArchiveTarBz2 = iota
	ArchiveZip    = iota
)

func ArchiveTypeName(archive_type int) string {
	switch archive_type {
	case ArchiveTarBz2:
		return "tar-bz2"
	case ArchiveZip:
		return "zip"
	default:
		return "UNKNOWN"
	}
}
