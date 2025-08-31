package database

import "github.com/OtchereDev/ris-common-sdk/pkg/io-dicom/media"

type Database interface {
	AddPatient(dcmObj media.DcmObj) error
	AddStudy(dcmObj media.DcmObj) error
	AddSeries(dcmObj media.DcmObj) error
	AddInstance(dcmObj media.DcmObj) error
	AddDicom(dcmObj media.DcmObj) error
}
