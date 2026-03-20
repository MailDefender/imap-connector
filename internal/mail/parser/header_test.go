package parser

import (
	"log"
	"testing"
)

func TestCleanHeader(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "This is a header with \r carriage return",
			expected: "This is a header with carriage return",
		},
		{
			input:    "This is a header with \t tab",
			expected: "This is a header with tab",
		},
		{
			input:    "This  is  a  header  with  multiple  spaces",
			expected: "This is a header with multiple spaces",
		},
		{
			input:    "This  is  a  header  with  multiple  spaces, tabs\t and\r carriage returns",
			expected: "This is a header with multiple spaces, tabs and carriage returns",
		},
	}

	for _, test := range tests {
		out := cleanHeader(test.input)
		if out != test.expected {
			t.Errorf("cleanHeader(%q) = %q, want %q", test.input, out, test.expected)
		}
	}
}

func TestParseHeader(t *testing.T) {
	input := `Return-Path: <test@sample.raw>
Delivered-To: deliver@sample.raw
Received: from localhost (HELO queue) (127.0.0.1)
	by localhost with SMTP; 21 Dec 2025 19:35:31 +0200
Received: from unknown (HELO output31.mail.sample.com) (127.0.0.1)
  by 127.0.0.1 with AES256-GCM-SHA384 encrypted SMTP; 21 Dec 2025 19:35:31 +0200
Received: from sample.mail.sample.com (unknown [127.0.0.1])
	by out.sample.com (Postfix) with ESMTP id abc123
	for <recipient@example.com>; Sun, 21 Dec 2025 17:35:31 +0000 (UTC)
Received: from in.mail.example.com (unknown [127.0.0.1])
	by vr.mail.example.com (Postfix) with ESMTP id abc1234
	for <recipient@example.com>; Sun, 21 Dec 2025 17:35:31 +0000 (UTC)
Received: from mail-out.example.com (mail-out.example.com [127.0.0.1])
	by in.example.com (Postfix) with ESMTPS id abc123
	for <recipient@example.com>; Sun, 21 Dec 2025 17:35:31 +0000 (UTC)
Authentication-Results: mail.example.com;
    arc=none (no signatures found);
    dkim=pass (2048-bit rsa key sha256) header.d=example.com header.i=@example.com header.b=abc123 header.a=rsa-sha256 header.s=random-selector;
    dmarc=none policy.published-domain-policy=none policy.applied-disposition=none policy.evaluated-disposition=none (p=none,d=none,d.eval=none) policy.policy-from=p header.from=example.fr;
    spf=pass smtp.mailfrom=me@example.com smtp.helo=mail-out.example.com;
    x-tls=pass smtp.version=TLSv1.3 smtp.cipher=TLS_AES_256_GCM_SHA384 smtp.bits=256
Received-SPF: pass
    (example.com: Sender is authorized to use 'me@example.com' in 'mfrom' identity (mechanism 'include:mail.example.com' matched))
    receiver=in77.mail.example.com;
    identity=mailfrom;
    envelope-from="me@example.com";
    helo=13.mail-out.example.com;
    client-ip=127.0.0.1
Received: from mail-out.example.com (unknown [127.0.0.1])
	by mail-out.example.com (Postfix) with ESMTP id abc123
`

	expectedHeaders := map[string][]string{
		"Return-Path":  {"<test@sample.raw>"},
		"Delivered-To": {"deliver@sample.raw"},
		"Received": {
			`from localhost (HELO queue) (127.0.0.1) by localhost with SMTP; 21 Dec 2025 19:35:31 +0200`,
			`from unknown (HELO output31.mail.sample.com) (127.0.0.1) by 127.0.0.1 with AES256-GCM-SHA384 encrypted SMTP; 21 Dec 2025 19:35:31 +0200`,
			`from sample.mail.sample.com (unknown [127.0.0.1]) by out.sample.com (Postfix) with ESMTP id abc123 for <recipient@example.com>; Sun, 21 Dec 2025 17:35:31 +0000 (UTC)`,
			`from in.mail.example.com (unknown [127.0.0.1]) by vr.mail.example.com (Postfix) with ESMTP id abc1234 for <recipient@example.com>; Sun, 21 Dec 2025 17:35:31 +0000 (UTC)`,
			`from mail-out.example.com (mail-out.example.com [127.0.0.1]) by in.example.com (Postfix) with ESMTPS id abc123 for <recipient@example.com>; Sun, 21 Dec 2025 17:35:31 +0000 (UTC)`,
			`from mail-out.example.com (unknown [127.0.0.1]) by mail-out.example.com (Postfix) with ESMTP id abc123`,
		},
		"Authentication-Results": {`mail.example.com; arc=none (no signatures found); dkim=pass (2048-bit rsa key sha256) header.d=example.com header.i=@example.com header.b=abc123 header.a=rsa-sha256 header.s=random-selector; dmarc=none policy.published-domain-policy=none policy.applied-disposition=none policy.evaluated-disposition=none (p=none,d=none,d.eval=none) policy.policy-from=p header.from=example.fr; spf=pass smtp.mailfrom=me@example.com smtp.helo=mail-out.example.com; x-tls=pass smtp.version=TLSv1.3 smtp.cipher=TLS_AES_256_GCM_SHA384 smtp.bits=256`},
		"Received-SPF": {
			`pass (example.com: Sender is authorized to use 'me@example.com' in 'mfrom' identity (mechanism 'include:mail.example.com' matched)) receiver=in77.mail.example.com; identity=mailfrom; envelope-from="me@example.com"; helo=13.mail-out.example.com; client-ip=127.0.0.1`,
		},
	}

	parser := NewHeaderParser(input)
	if err := parser.Parse(); err != nil {
		log.Fatal("an error occured while parsing headers")
	}

	for key, val := range expectedHeaders {
		get := parser.Get(key)

		if len(val) != len(get) {
			log.Fatalf("the expected length of the string array headers does not match for %s", key)
		}

		for i := 0; i < len(val); i++ {
			if get[i] != val[i] {
				log.Fatalf("the retrieved header value does not match with the expected one ; %s != %s", get[i], val[i])
			}
		}
	}
}
