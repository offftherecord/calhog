package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Attendee struct {
	ID               string `json:"id"`
	Email            string `json:"email"`
	Displayname      string `json:"displayName"`
	Organizer        bool   `json:"organizer"`
	Self             bool   `json:"self"`
	Resource         bool   `json:"resource"`
	Optional         bool   `json:"optional"`
	Responsestatus   string `json:"responseStatus"`
	Comment          string `json:"comment"`
	Additionalguests int    `json:"additionalGuests"`
}

type Attachment struct {
	Fileurl  string `json:"fileUrl"`
	Title    string `json:"title"`
	Mimetype string `json:"mimeType"`
	Iconlink string `json:"iconLink"`
	Fileid   string `json:"fileId"`
}

type Creator struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Displayname string `json:"displayName"`
	Self        bool   `json:"self"`
}
type Organizer struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Displayname string `json:"displayName"`
	Self        bool   `json:"self"`
}

type Start struct {
	Date     string `json:"date"`
	Datetime string `json:"dateTime"`
	Timezone string `json:"timeZone"`
}

type End struct {
	Date     string `json:"date"`
	Datetime string `json:"dateTime"`
	Timezone string `json:"timeZone"`
}

type Originalstarttime struct {
	Date     string `json:"date"`
	Datetime string `json:"dateTime"`
	Timezone string `json:"timeZone"`
}

type Gadget struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Link     string `json:"link"`
	Iconlink string `json:"iconLink"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Display  string `json:"display"`
}

type Overrides struct {
	Method  string `json:"method"`
	Minutes int    `json:"minutes"`
}

type Reminders struct {
	Usedefault bool        `json:"useDefault"`
	Overrides  []Overrides `json:"overrides"`
}

type Source struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type Defaultreminders struct {
	Method  string `json:"method"`
	Minutes int    `json:"minutes"`
}
type Response struct {
	Kind             string             `json:"kind"`
	Etag             string             `json:"etag"`
	Summary          string             `json:"summary"`
	Description      string             `json:"description"`
	Updated          string             `json:"updated"`
	Timezone         string             `json:"timeZone"`
	Accessrole       string             `json:"accessRole"`
	Defaultreminders []Defaultreminders `json:"defaultReminders"`
	Nextpagetoken    string             `json:"nextPageToken"`
	Nextsynctoken    string             `json:"nextSyncToken"`
	Items            []Item             `json:"items"`
}

type Key struct {
	Type string `json:"type"`
}

type Conferencesolution struct {
	Key     Key    `json:"key"`
	Name    string `json:"name"`
	Iconuri string `json:"iconUri"`
}

type Entrypoint struct {
	Entrypointtype string `json:"entryPointType"`
	URI            string `json:"uri"`
	Label          string `json:"label"`
	Pin            string `json:"pin"`
	Accesscode     string `json:"accessCode"`
	Meetingcode    string `json:"meetingCode"`
	Passcode       string `json:"passcode"`
	Password       string `json:"password"`
}

type Status struct {
	Statuscode string `json:"statusCode"`
}

type Conferencesolutionkey struct {
	Type string `json:"type"`
}

type Createrequest struct {
	Requestid             string                `json:"requestId"`
	Conferencesolutionkey Conferencesolutionkey `json:"conferenceSolutionKey"`
	Status                Status                `json:"status"`
}

type Conferencedata struct {
	Createrequest      Createrequest      `json:"createRequest"`
	Entrypoints        []Entrypoint       `json:"entryPoints"`
	Conferencesolution Conferencesolution `json:"conferenceSolution"`
	Conferenceid       string             `json:"conferenceId"`
	Signature          string             `json:"signature"`
	Notes              string             `json:"notes"`
}

type Item struct {
	Kind                    string            `json:"kind"`
	Etag                    string            `json:"etag"`
	ID                      string            `json:"id"`
	Status                  string            `json:"status"`
	Htmllink                string            `json:"htmlLink"`
	Created                 string            `json:"created"`
	Updated                 string            `json:"updated"`
	Summary                 string            `json:"summary"`
	Description             string            `json:"description"`
	Location                string            `json:"location"`
	Colorid                 string            `json:"colorId"`
	Creator                 Creator           `json:"creator"`
	Organizer               Organizer         `json:"organizer"`
	Start                   Start             `json:"start"`
	End                     End               `json:"end"`
	Endtimeunspecified      bool              `json:"endTimeUnspecified"`
	Recurrence              []string          `json:"recurrence"`
	Recurringeventid        string            `json:"recurringEventId"`
	Originalstarttime       Originalstarttime `json:"originalStartTime"`
	Transparency            string            `json:"transparency"`
	Visibility              string            `json:"visibility"`
	Icaluid                 string            `json:"iCalUID"`
	Sequence                int               `json:"sequence"`
	Attendees               []Attendee        `json:"attendees"`
	Attendeesomitted        bool              `json:"attendeesOmitted"`
	Hangoutlink             string            `json:"hangoutLink"`
	Conferencedata          Conferencedata    `json:"conferenceData"`
	Gadget                  Gadget            `json:"gadget"`
	Anyonecanaddself        bool              `json:"anyoneCanAddSelf"`
	Guestscaninviteothers   bool              `json:"guestsCanInviteOthers"`
	Guestscanmodify         bool              `json:"guestsCanModify"`
	Guestscanseeotherguests bool              `json:"guestsCanSeeOtherGuests"`
	Privatecopy             bool              `json:"privateCopy"`
	Locked                  bool              `json:"locked"`
	Reminders               Reminders         `json:"reminders"`
	Source                  Source            `json:"source"`
	Attachments             []Attachment      `json:"attachments"`
	Eventtype               string            `json:"eventType"`
}

type Finding struct {
	HtmlLink       string
	Type           string
	User           string
	Email          string
	ZoomLink       string
	HangoutLink    string
	AttachmentLink string
	Summary        string
	Date           string
	Title          string
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func searchEmails(item Item, emails []string) []string {
	if item.Creator.Email != "" && !contains(emails, item.Creator.Email) {
		emails = append(emails, item.Creator.Email)
		report(Finding{Type: "email", HtmlLink: item.Htmllink, Date: item.Start.Datetime, Summary: item.Summary, Email: item.Creator.Email})
	}
	if item.Organizer.Email != "" && !contains(emails, item.Organizer.Email) {
		emails = append(emails, item.Organizer.Email)
		report(Finding{Type: "email", HtmlLink: item.Htmllink, Date: item.Start.Datetime, Summary: item.Summary, Email: item.Organizer.Email})
	}
	if len(item.Attendees) > 0 {
		for _, i := range item.Attendees {
			if !contains(emails, i.Email) {
				emails = append(emails, i.Email)
				report(Finding{Type: "email", HtmlLink: item.Htmllink, Date: item.Start.Datetime, Summary: item.Summary, Email: i.Email})
			}
		}
	}

	return emails
}

func searchHangouts(item Item, hangouts []string) []string {
	if item.Hangoutlink != "" && !contains(hangouts, item.Hangoutlink) {
		hangouts = append(hangouts, item.Hangoutlink)
		report(Finding{Type: "hangout", HtmlLink: item.Htmllink, Date: item.Start.Datetime, Summary: item.Summary, HangoutLink: item.Hangoutlink})

	}

	return hangouts
}

func searchAttachments(item Item, attachments []string) []string {
	if len(item.Attachments) > 0 {
		for _, attachment := range item.Attachments {
			attachments = append(attachments, attachment.Fileurl)
			report(Finding{Type: "attachment", HtmlLink: item.Htmllink, Summary: item.Summary, AttachmentLink: attachment.Fileurl, Title: attachment.Title, Date: item.Start.Datetime})

		}
	}
	return attachments
}

func searchZoomlinks(item Item, zoomlinks []string) []string {
	if item.Location != "" && strings.Contains(item.Location, "zoom") && !contains(zoomlinks, item.Location) {
		zoomlinks = append(zoomlinks, item.Location)
		report(Finding{Type: "zoom", ZoomLink: item.Location, Date: item.Start.Datetime, Summary: item.Summary})
	}

	return zoomlinks
}

func request(url string) Response {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var res Response
	json.Unmarshal(b, &res)

	return res
}

func report(f Finding) {
	l, err := json.Marshal(&f)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(l))
}

func main() {
	key := "AIzaSyBNlYH01_9Hc5S1J9vuFmu2nUqBZJNAXxs" // Auto generated from Google
	pageToken := ""
	gcalAPI := "https://www.googleapis.com/calendar/v3/calendars/%s/events?singleEvents=true&&maxAttendees=1000&maxResults=2500&sanitizeHtml=true&key=%s"
	var emails []string
	var hangouts []string
	var attachments []string
	var zoomlinks []string

	scanner := bufio.NewScanner(os.Stdin)

	if scanner.Err() != nil {
		log.Panic(scanner.Err())
	}

	for scanner.Scan() {
		user := scanner.Text()
		for {
			var url string
			if pageToken == "" {
				url = fmt.Sprintf(gcalAPI, user, key)
			} else {
				url = fmt.Sprintf(gcalAPI, user, key) + fmt.Sprintf("&pageToken=%s", pageToken)
			}

			res := request(url)
			pageToken = res.Nextpagetoken

			for _, item := range res.Items {
				emails = searchEmails(item, emails)
				hangouts = searchHangouts(item, hangouts)
				attachments = searchAttachments(item, attachments)
				zoomlinks = searchZoomlinks(item, zoomlinks)
			}

			// End of calendar pagination
			if res.Nextsynctoken != "" {
				break
			}
		}
	}

}
