package create

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/umbrella-sh/um-common/jsons"

	"github.com/umbrella-sh/simply-dns-cli/internal/api"
	apio "github.com/umbrella-sh/simply-dns-cli/internal/api_objects"
	"github.com/umbrella-sh/simply-dns-cli/internal/collectors"
	"github.com/umbrella-sh/simply-dns-cli/internal/forms"
	"github.com/umbrella-sh/simply-dns-cli/internal/shared"
	"github.com/umbrella-sh/simply-dns-cli/internal/styles"
)

func cmdRun(_ *cobra.Command, _ []string) {
	styles.Println(styles.Info("Add new dns record"))
	styles.Blank()

	cancelled, domain, record := collectInfo()
	if cancelled {
		printCancelText()
		return
	}

	styles.Blank()
	var accepted bool
	cancelled, accepted = collectors.AcceptInfo()
	if !accepted {
		printNotAcceptedText()
		return
	}
	styles.Blank()

	createRecord(domain, record)
}

//goland:noinspection GoNameStartsWithPackageName
func createRecord(domain string, record *apio.SimplyDnsRecord) {
	styles.WaitPrint("Creating dns entry")

	res, err := api.CreateDnsRecord(domain, record)
	if err != nil {
		styles.FailPrint("Failed to create DNS Entry")
		styles.FailPrint("Error: %v", err)
		return
	}

	if res.Status != 200 {
		styles.FailPrint("Failed to create DNS Entry")
		styles.FailPrint("Error: %d, %v", res.Status, res.Message)
		return
	}

	styles.SuccessPrint("DNS Entry created on %s", domain)
}

func collectInfo() (cancelled bool, domain string, record *apio.SimplyDnsRecord) {
	record = &apio.SimplyDnsRecord{}

	cancelled, domain = collectors.CollectDomain(options.Domain)
	if cancelled {
		return cancelled, "", nil
	}

	if options.Type == "" {
		cancelled, record.Type = forms.RunTypeSelect("")
		if cancelled {
			return cancelled, "", nil
		}
	} else {
		record.Type = apio.DnsRecordType(options.Type)
		shared.PrintValue(forms.TypeSelectHeader, apio.DnsTypeToText(record.Type))
	}

	if options.TTL <= 0 {
		cancelled, record.TTL = forms.RunTtlSelect(apio.DnsRecTTLHour1)
		if cancelled {
			return cancelled, "", nil
		}
	} else {
		record.TTL = apio.DnsRecordTTL(options.TTL)
		shared.PrintValue(forms.TtlSelectHeader, apio.DnsTTLToText(record.TTL))
	}

	if options.Name == "" {
		var name string
		cancelled, name = forms.RunNameInput("")
		if cancelled {
			return cancelled, "", nil
		}
		record.Name = fmt.Sprintf("%s.%s", name, domain)
	} else {
		record.Name = fmt.Sprintf("%s.%s", options.Name, domain)
		shared.PrintValue(forms.NameInputHeader, record.Name)
	}

	if options.Data == "" {
		cancelled, record.Data = forms.RunDataInput("")
		if cancelled {
			return cancelled, "", nil
		}
	} else {
		record.Data = options.Data
		shared.PrintValue(forms.DataInputHeader, record.Data)
	}

	if record.Type == apio.DnsRecTypeMX {
		if options.Priority <= 0 {
			cancelled, record.Priority = forms.RunPriorityInput(nil)
			if cancelled {
				return cancelled, "", nil
			}
		} else {
			record.Priority = jsons.NewJsonInt32(int32(options.Priority))
			shared.PrintValue(forms.PriorityInputHeader, record.Priority.ToString())
		}
	} else {
		record.Priority = jsons.NullJsonInt32()
	}

	if options.Comment == NoCommentValue {
		cancelled, record.Comment = forms.RunCommentInput("")
		if cancelled {
			return cancelled, "", nil
		}
	} else {
		record.Comment = options.Comment
		shared.PrintValue(forms.CommentInputHeader, record.Comment)
	}

	return cancelled, domain, record
}

func printCancelText() { styles.Println(styles.Warn("\nCreate was cancelled\n")) }
func printNotAcceptedText() {
	styles.Println(styles.Warn("\nInformation is not accepted and no record was created\n"))
}
