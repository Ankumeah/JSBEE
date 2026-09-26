package emails

import (
	"bytes"
	"fmt"
)

func PaperAcceptedEmail(
	authorName string,
	paperTitle string,
) (HTML string, text string, err error) {
	buf := bytes.NewBuffer(nil)
	if err := templates[PaperAcceptedFile].Execute(
		buf,
		struct {
			AuthorName      string
			PaperTitle      string
			SiteName        string
			SiteDescription string
			HomeLink        string
			VolumesLink     string
		}{
			AuthorName:      authorName,
			PaperTitle:      paperTitle,
			SiteName:        siteName,
			SiteDescription: siteDescription,
			HomeLink:        homeLink,
			VolumesLink:     volumesLink,
		},
	); err != nil {
		return "", "", err
	}

	HTML = buf.String()
	text = fmt.Sprintf(
		"Congratulations %s! your paper was accepted.\n"+
			"\n"+
			"Our editors have reviewed your paper (%s) and approved its contents."+
			"Thank you for contributing your research to "+siteName+".\n"+
			"\n"+
			"What happens next\n"+
			"1. Your paper waits in the queue for the next publish issue.\n"+
			"2. Your paper receives a number, volume and an issue.\n"+
			"3. You receive another email as soon as it is published.\n"+
			"\n"+
			"Warm regards,\n"+
			" - The "+siteName+" editorial team\n",
		authorName, paperTitle,
	)

	return
}

func PaperRejectedEmail(
	authorName string,
	paperTitle string,
	reason string,
) (HTML string, text string, err error) {
	buf := bytes.NewBuffer(nil)
	if err := templates[PaperRejectedFile].Execute(
		buf,
		struct {
			AuthorName      string
			PaperTitle      string
			Reason          string
			SiteName        string
			SiteDescription string
			HomeLink        string
			VolumesLink     string
		}{
			AuthorName:      authorName,
			PaperTitle:      paperTitle,
			Reason:          reason,
			SiteName:        siteName,
			SiteDescription: siteDescription,
			HomeLink:        homeLink,
			VolumesLink:     volumesLink,
		},
	); err != nil {
		return "", "", err
	}

	HTML = buf.String()
	text = fmt.Sprintf(
		"An update on your submission at "+siteName+"\n"+
			"\n"+
			"%s\n"+
			"Thank you for submitting your work to "+siteName+"."+
			"After careful review, our editors have decided unfortunately not to move this paper forward to publication.\n"+
			"%s"+
			"\n"+
			"This decision reflects fit and quality standards for the journal, not the value of your research."+
			"You are welcome to revise and submit again, or send a new paper in the future.\n"+
			"\n"+
			"Warm regards,\n"+
			" - The "+siteName+" editorial team\n",
		paperTitle, feedbackLine(reason),
	)

	return
}

func PaperPublishedEmail(
	authorName string,
	paperTitle string,
) (HTML string, text string, err error) {
	buf := bytes.NewBuffer(nil)
	if err := templates[PaperPublishedFile].Execute(
		buf,
		struct {
			AuthorName      string
			PaperTitle      string
			SiteName        string
			SiteDescription string
			HomeLink        string
			VolumesLink     string
		}{
			AuthorName:      authorName,
			PaperTitle:      paperTitle,
			SiteName:        siteName,
			SiteDescription: siteDescription,
			HomeLink:        homeLink,
			VolumesLink:     volumesLink,
		},
	); err != nil {
		return "", "", err
	}

	HTML = buf.String()
	text = fmt.Sprintf(
		"Congratulations %s! Your paper is live.\n"+
			"\n"+
			"Your paper (%s) has been published in "+siteName+" and is now free for anyone to read and download."+
			"Thank you for sharing your research with the community.\n"+
			"\n"+
			"View your paper at by visiting "+volumesLink+
			"\n"+
			"Warm regards,\n"+
			" - The "+siteName+" editorial team\n",
		authorName, paperTitle,
	)

	return
}

func feedbackLine(reason string) string {
	if reason == "" {
		return ""
	}
	return "\nEditor feedback: " + reason + "\n"
}
