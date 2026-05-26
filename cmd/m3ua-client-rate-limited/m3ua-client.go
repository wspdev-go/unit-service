package main

import (
	"context"
	"encoding/hex"
	"flag"
	"io"
	"sync"
	"time"

	"gitlab.horisen.org/horisen/ss7/hstplib/logger"
	"gitlab.horisen.org/horisen/ss7/hstplib/thirdparty/wmnsk/go-m3ua"
	"gitlab.horisen.org/horisen/ss7/hstplib/thirdparty/wmnsk/go-m3ua/messages/params"
	"gitlab.horisen.org/horisen/ss7/hstplib/thirdparty/wmnsk/go-sccp"
	sccpParams "gitlab.horisen.org/horisen/ss7/hstplib/thirdparty/wmnsk/go-sccp/params"
	"gitlab.horisen.org/horisen/ss7/hstplib/thirdparty/wmnsk/go-sccp/utils"
	"gitlab.horisen.org/horisen/ss7/hstplib/thirdparty/wmnsk/go-tcap"

	"github.com/ishidawataru/sctp"
)

func main() {
	var (
		localAddr = flag.String("laddr", "127.0.0.1:2915", "Local IP and Port to connect to.")
		addr      = flag.String("addr", "127.0.0.1:2905", "Remote IP and Port to connect to.")
		opc       = flag.Uint("opc", 4369, "Originating point code.")
		dpc       = flag.Uint("dpc", 546, "Destination point code.") // Server OPC
		//dpc  = flag.Uint("dpc", 546, "Destination point code.") // HSTP OPC

		// one-to-one to server
		//addr = flag.String("addr", "127.0.0.1:2906", "Remote IP and Port to connect to.")
		//opc  = flag.Uint("opc", 546, "Originating point code.")
		//dpc  = flag.Uint("dpc", 273, "Destination point code.")

		//data  = flag.String("data", "098003101b0d120800110488220880000000010b1208001204947101670700cd6281ca48030123046c81c2a181bf02010002012e3081b6800862428345350073f28407919471016700000481a0440d91945171447709f4000052102051630240a0050003050201906136fb1d5224c768d07c8c0e9bcd6550790e4297eb7490333c46b7d3747af80c22bfc768903b3d46d35da0a29b3e1fa3eb6c72fa5c779f5d2045728c06ddfd727219d47ecbcf6537681c76816233dd0cf68ad1743058152d07adeb723dc8fe968bcbe9f5bbdd2ebb41753719744fcb41edf0185d7683ca6977d90da296e5edb41b64f6cb41", "Payload to send on M3UA in hex stream format.")
		data   = flag.String("data", "aaaaaaaa", "Payload to send on M3UA in hex stream format.")
		hbInt  = flag.Duration("hb-interval", 0, "Interval for M3UA BEAT. Put 0 to disable")
		n      = flag.Int("n", 1, "Number of connections.")
		r      = flag.Int("r", 8000, "Rate limit per second.")
		cdPaGT = flag.String("cdpa-gt", "123456789012345", "Called party global title.")
		cgPaGT = flag.String("cgpa-gt", "100000001", "Calling party global title.")
		otid   = flag.Int("otid", 123456, "Originating transaction ID.")

		otidFrom = flag.Int("otid-from", 0, "Originating transaction ID from.")
		otidTo   = flag.Int("otid-to", 0, "Originating transaction ID to.")

		gti = flag.Uint("gti", uint(sccpParams.GTITTNPESNAI), "Global Title Indicator (GTI) for Called Party Address (CdPA).")
		npi = flag.Uint("npi", uint(sccpParams.NPISDNTelephony), "Numbering Plan Indicator (NPI) for Called Party Address (CdPA).")
		nai = flag.Uint("nai", uint(sccpParams.NAIInternationalNumber), "Numbering Address Indicator (NAI) for Called Party Address (CdPA).")
		tt  = flag.Uint("tt", 0, "Translation Type (TT) for Called Party Address (CdPA).")
	)
	flag.Parse()

	if *otidFrom == 0 && *otidTo == 0 {
		// if otidFrom and otidTo are not set, use the single otid value
		otidFrom = otid
		otidTo = otid
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; i < *n; i++ {
		go serve(ctx, *localAddr, *addr, *opc, *dpc, *data, *hbInt, *r, *cdPaGT, *cgPaGT, *otid, *otidFrom, *otidTo, *gti, *npi, *nai, *tt)
	}

	<-ctx.Done()
}

func createSCCPMessage(cdPaGT string, cgPaGT string, otid int, gti, npi, nai, tt uint) ([]byte, error) {
	payload, err := hex.DecodeString("040800010121436587f9")
	if err != nil {
		return nil, err
	}

	tcapBytes, err := tcap.NewBeginInvokeWithDialogue(
		//uint32(0x11111111),               // Originating Transaction ID
		uint32(otid),                     // Originating Transaction ID
		tcap.DialogueAsID,                // DialogueType
		tcap.LocationCancellationContext, // ACN
		3,                                // ACN Version
		0,                                // Invoke Id
		3,                                // OpCode
		payload,                          // TCAP payload
	).MarshalBinary()

	cdGti := sccpParams.GlobalTitleIndicator(gti)
	cdAi := sccpParams.NewAddressIndicator(false, true, false, cdGti)
	cdPA := sccpParams.NewCalledPartyAddress(
		cdAi, 0, 6, sccpParams.NewGlobalTitle(
			cdGti,
			sccpParams.TranslationType(tt),
			sccpParams.NumberingPlan(npi),
			sccpParams.ESBCDOdd,
			sccpParams.NatureOfAddressIndicator(nai),
			utils.MustBCDEncode(cdPaGT),
		),
	)

	cgGti := sccpParams.GTITTNPESNAI
	cgAi := sccpParams.NewAddressIndicator(false, true, false, cgGti)
	cgPA := sccpParams.NewCallingPartyAddress(
		cgAi, 0, 7, sccpParams.NewGlobalTitle(
			cgGti,
			sccpParams.TranslationType(1),
			sccpParams.NPISDNMobile,
			sccpParams.ESBCDOdd,
			sccpParams.NAIInternationalNumber,
			utils.MustBCDEncode(cgPaGT),
		),
	)
	// create UDT message with CdPA, CgPA and payload
	udt := sccp.NewUDT(
		1,    // Protocol Class
		true, // Message handling
		cdPA,
		cgPA,
		tcapBytes, // SCCP payload
	)

	u, err := udt.MarshalBinary()

	return u, err
}

var mu sync.Mutex

func serve(pCtx context.Context, localAddr string, addr string, opc uint, dpc uint, data string, hbInt time.Duration, rate int,
	cdPaGT string, cgPaGT string, otid int, otidFrom int, otidTo int, gti, npi, nai, tt uint) {

	ctx, cancel := context.WithCancel(pCtx)
	defer cancel()

	// create *Config to be used in M3UA connection
	config := m3ua.NewConfig(
		uint32(opc),           // OriginatingPointCode
		uint32(dpc),           // DestinationPointCode
		params.ServiceIndSCCP, // ServiceIndicator
		0,                     // NetworkIndicator
		0,                     // MessagePriority
		1,                     // SignalingLinkSelection
	)
	config. // set parameters to use
		EnableHeartbeat(hbInt, 10*time.Second).
		SetAspIdentifier(1).
		SetTrafficModeType(params.TrafficModeLoadshare).
		SetNetworkAppearance(0).
		SetRoutingContexts(1, 2)

	// setup SCTP peer on the specified IPs and Port.

	raddr, err := sctp.ResolveSCTPAddr("sctp", addr)
	if err != nil {
		logger.Error("Failed to resolve SCTP address: %s", err)
		return
	}

	laddr, err := sctp.ResolveSCTPAddr("sctp", localAddr)
	if err != nil {
		logger.Error("Failed to resolve SCTP address: %s", err)
		return
	}

	var conn *m3ua.Conn
	defer func() {
		if conn == nil {
			return
		}
		err = conn.Close()
		if err != nil {
			logger.Error("Failed to close connection: %s", err)
		}
	}()

	currentTcapId := otidFrom
	getTcapId := func() int {
		mu.Lock()
		if currentTcapId > otidTo {
			currentTcapId = otidFrom
		}
		tcapId := currentTcapId
		currentTcapId++
		mu.Unlock()
		return tcapId
	}

	//maxConnectionAttempts := 1000
	//currentConnectionAttempt := 0
	for {
		//currentConnectionAttempt++

		conn, err = m3ua.Dial(ctx, "m3ua", laddr, raddr, config)
		if err != nil {
			//logger.Error("Failed to dial M3UA: %s (%d attempt)", err, currentConnectionAttempt)
			logger.Error("Failed to dial M3UA: %s", err)
			time.Sleep(5 * time.Second)
			continue
		} else {
			logger.Info("Connected")
			//break
		}

		if rate >= 100 {
			for i := 0; i < rate/100; i++ {
				go write100perSecond(conn, cdPaGT, cgPaGT, otidFrom, otidTo, gti, npi, nai, tt, getTcapId)
			}
		} else if rate >= 10 {
			for i := 0; i < rate/10; i++ {
				go write10perSecond(conn, cdPaGT, cgPaGT, otidFrom, otidTo, gti, npi, nai, tt, getTcapId)
			}
		} else {
			for i := 0; i < rate; i++ {
				go write1perSecond(conn, cdPaGT, cgPaGT, otidFrom, otidTo, gti, npi, nai, tt, getTcapId)
			}
		}

		read(conn)

	}

	// For slow rate
	//for {
	//	if _, err = conn.Write(payload); err != nil {
	//		logger.Error("Failed to write M3UA data: %s", err)
	//	}
	//	ClientSendTotalVec.WithLabelValues("client_write").Inc()
	//	time.Sleep(1 * time.Second)
	//}

	<-ctx.Done()
}

func write100perSecond(conn *m3ua.Conn, cdPaGT string, cgPaGT string, otidFrom int, otidTo int, gti, npi, nai, tt uint, getTcapId func() int) {

	writeToConn(conn, 10*time.Millisecond, 5000, cdPaGT, cgPaGT, gti, npi, nai, tt, getTcapId)
}

func write10perSecond(conn *m3ua.Conn, cdPaGT string, cgPaGT string, otidFrom int, otidTo int, gti, npi, nai, tt uint, getTcapId func() int) {

	writeToConn(conn, 100*time.Millisecond, 500, cdPaGT, cgPaGT, gti, npi, nai, tt, getTcapId)
}

func write1perSecond(conn *m3ua.Conn, cdPaGT string, cgPaGT string, otidFrom int, otidTo int, gti, npi, nai, tt uint, getTcapId func() int) {

	writeToConn(conn, 1*time.Second, 50, cdPaGT, cgPaGT, gti, npi, nai, tt, getTcapId)
}

func writeToConn(conn *m3ua.Conn, timeDuration time.Duration, counter int, cdPaGT string, cgPaGT string, gti, npi, nai, tt uint, getTcapId func() int) {
	limiter := time.Tick(timeDuration)
	for range counter {
		otid := getTcapId()
		d, err := createSCCPMessage(cdPaGT, cgPaGT, otid, gti, npi, nai, tt)
		if err != nil {
			logger.Error("Failed to createSCCPMessage: %s", err)
		}
		<-limiter
		if _, err = conn.WriteToStream(d, m3ua.DataStreamId); err != nil {
			logger.Error("Failed to write M3UA data: %s", err)
			return
		}
		ClientSendTotalVec.WithLabelValues("client_write").Inc()
		ClientSendTotal.Inc()
	}

	logger.Info("finished %d per second, send %d messages, time %d seconds", timeDuration, counter, timeDuration)
}

func read(conn *m3ua.Conn) {
	buf := make([]byte, 4096)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			// this indicates the connection is no longer alive. close M3UA conn and wait for INIT again.
			if err == io.EOF {
				logger.Warn("Closed M3UA conn with: %s, waiting to come back...", conn.RemoteAddr())
				return
			}
			// this indicates some unexpected error occurred on M3UA conn.
			logger.Error("Error reading from M3UA connection: %s", err)
			return
		}
		ClientReceiveTotalVec.WithLabelValues("client_read").Inc()
		ClientReceiveTotal.Inc()
	}
}
