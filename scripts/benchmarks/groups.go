package main

// Explicit workloads, friendly labels, and fixed comparison references.
// Add new benchmark cases here before publishing their documentation.
var groups = map[string][]groupDef{
	"cache": {
		{
			Title:       "Set and get by payload size",
			Description: "One operation sets a prepared value, gets it back, and checks the returned bytes over local container TCP connections.",
			Baseline:    "Memcached",
			Operation:   "set/get",
			Note:        "Allocated bytes belong to the Go client process; they exclude server memory. Redis persistence is disabled.",
			Cases: []caseDef{
				{Label: "2 KiB", Rows: []rowDef{
					{"Memcached", "cache", "BenchmarkMemcache2KB"},
					{"Redis", "cache", "BenchmarkRedis2KB"},
				}},
				{Label: "4 KiB", Rows: []rowDef{
					{"Memcached", "cache", "BenchmarkMemcache4KB"},
					{"Redis", "cache", "BenchmarkRedis4KB"},
				}},
				{Label: "8 KiB", Rows: []rowDef{
					{"Memcached", "cache", "BenchmarkMemcache8KB"},
					{"Redis", "cache", "BenchmarkRedis8KB"},
				}},
				{Label: "16 KiB", Rows: []rowDef{
					{"Memcached", "cache", "BenchmarkMemcache16KB"},
					{"Redis", "cache", "BenchmarkRedis16KB"},
				}},
				{Label: "32 KiB", Rows: []rowDef{
					{"Memcached", "cache", "BenchmarkMemcache32KB"},
					{"Redis", "cache", "BenchmarkRedis32KB"},
				}},
				{Label: "64 KiB", Rows: []rowDef{
					{"Memcached", "cache", "BenchmarkMemcache64KB"},
					{"Redis", "cache", "BenchmarkRedis64KB"},
				}},
				{Label: "128 KiB", Rows: []rowDef{
					{"Memcached", "cache", "BenchmarkMemcache128KB"},
					{"Redis", "cache", "BenchmarkRedis128KB"},
				}},
				{Label: "256 KiB", Rows: []rowDef{
					{"Memcached", "cache", "BenchmarkMemcache256KB"},
					{"Redis", "cache", "BenchmarkRedis256KB"},
				}},
				{Label: "512 KiB", Rows: []rowDef{
					{"Memcached", "cache", "BenchmarkMemcache512KB"},
					{"Redis", "cache", "BenchmarkRedis512KB"},
				}},
			},
		},
	},
	"compression": {
		{
			Title:       "Compress and decompress 1 MiB",
			Description: "One operation creates codec state and buffers, compresses a printable-ASCII input, then decompresses it. Inputs have the same size and distribution, but are generated separately for each codec.",
			Baseline:    "Gzip",
			Operation:   "round trip",
			Note:        "Gzip and Brotli use their best-compression settings; LZ4 uses level 9. Comparable compressed sizes were not recorded, so these timings do not establish a size-versus-speed winner.",
			Rows: []rowDef{
				{"Gzip", "compression", "BenchmarkCompressDecompressGzip"},
				{"Brotli", "compression", "BenchmarkCompressDecompressBrotli"},
				{"LZ4", "compression", "BenchmarkCompressDecompressLZ4"},
			},
		},
	},
	"crypto": {
		{
			Title:       "Signed JWT lifecycle",
			Description: "Create an HS256 JWT, or parse and verify a prepared JWT including its claims and expiry. Creation includes timestamp calls.",
			Operation:   "token",
			Note:        "Creation and verification perform different work. No equivalent implementation reference was measured.",
			Rows: []rowDef{
				{"Create and sign", "crypto/sign", "BenchmarkJWTGeneration"},
				{"Parse and verify", "crypto/sign", "BenchmarkJWTVerification"},
			},
		},
		{
			Title:       "Encrypted token lifecycle",
			Description: "Create an AES-256-GCM token with JSON encoding and a fresh nonce, or decrypt a prepared token and decode its payload and expiry.",
			Operation:   "token",
			Note:        "Encryption and JWT signing provide different guarantees. These lifecycle costs are kept separate from the JWT results.",
			Rows: []rowDef{
				{"Encode and encrypt", "crypto/encrypt", "BenchmarkGenerateEncryptedToken"},
				{"Decrypt and decode", "crypto/encrypt", "BenchmarkDecryptEncryptedToken"},
			},
		},
	},
	"db": {
		{
			Title:       "Insert a fixed value",
			Description: "Prepare and execute one insert of a fixed string; the table grows during timing.",
			Baseline:    "PostgreSQL",
			Operation:   "query",
			Note:        "Both clients use local container TCP connections. Server memory is excluded; configured durability and caching differ between engines.",
			Rows: []rowDef{
				{"PostgreSQL", "db/postgresql", "BenchmarkPostgresSingleInsertFixedData"},
				{"MariaDB", "db/mariadb", "BenchmarkMariaDBSingleInsertFixedData"},
			},
		},
		{
			Title:       "Generate and insert a UUID",
			Description: "Generate a UUID, then prepare and execute its insert; the table grows during timing.",
			Baseline:    "PostgreSQL",
			Operation:   "query",
			Note:        "Both clients use local container TCP connections. Server memory is excluded; configured durability and caching differ between engines.",
			Rows: []rowDef{
				{"PostgreSQL", "db/postgresql", "BenchmarkPostgresSingleInsertRandomData"},
				{"MariaDB", "db/mariadb", "BenchmarkMariaDBSingleInsertRandomData"},
			},
		},
		{
			Title:       "Indexed lookup in 1 million rows",
			Description: "Look up and verify one indexed value in a prepared one-million-row table. Population is excluded.",
			Baseline:    "PostgreSQL",
			Operation:   "query",
			Note:        "Both clients use local container TCP connections. Server memory is excluded; configured durability and caching differ between engines.",
			Rows: []rowDef{
				{"PostgreSQL", "db/postgresql", "BenchmarkPostgresInsert1MilAndFindMiddle"},
				{"MariaDB", "db/mariadb", "BenchmarkMariaDBInsert1MilAndFindMiddle"},
			},
		},
		{
			Title:       "Read and order 100,000 rows",
			Description: "Read 100,000 rows ordered by value, scan them, and verify their ordering and count.",
			Baseline:    "PostgreSQL",
			Operation:   "query",
			Note:        "Both clients use local container TCP connections. Server memory is excluded; configured durability and caching differ between engines.",
			Rows: []rowDef{
				{"PostgreSQL", "db/postgresql", "BenchmarkPostgresInsert100kAndSort"},
				{"MariaDB", "db/mariadb", "BenchmarkMariaDBInsert100kAndSort"},
			},
		},
		{
			Title:       "Filter and order 80,000 matches",
			Description: "Filter a prepared 100,000-row table by a fixed timestamp cutoff, scan the 80,000 matching rows, and verify their order and count.",
			Baseline:    "PostgreSQL",
			Operation:   "query",
			Note:        "Both clients use local container TCP connections. Server memory is excluded; configured durability and caching differ between engines.",
			Rows: []rowDef{
				{"PostgreSQL", "db/postgresql", "BenchmarkPostgresInsert100kAndQueryOlderThan15Min"},
				{"MariaDB", "db/mariadb", "BenchmarkMariaDBInsert100kAndQueryOlderThan15Min"},
			},
		},
		{
			Title:       "Embedded SQLite costs",
			Description: "The same insert and read workloads use an embedded engine and a private persistent database file. Read fixture population is excluded; scans include count and ordering checks.",
			Operation:   "query",
			Note:        "SQLite has no client/server TCP round trip. Its costs are shown separately from the networked database comparison; the five rows also perform different operations.",
			Rows: []rowDef{
				{"Insert a fixed value", "db/sqlite", "BenchmarkSQLiteSingleInsertFixedData"},
				{"Generate and insert a UUID", "db/sqlite", "BenchmarkSQLiteSingleInsertRandomData"},
				{"Indexed lookup in 1 million rows", "db/sqlite", "BenchmarkSQLiteInsert1MilAndFindMiddle"},
				{"Read and order 100,000 rows", "db/sqlite", "BenchmarkSQLiteInsert100kAndSort"},
				{"Filter and order 80,000 matches", "db/sqlite", "BenchmarkSQLiteInsert100kAndQueryOlderThan15Min"},
			},
		},
	},
	"encoding": {
		{
			Title:       "Encode by input size",
			Description: "One operation allocates an output buffer and encodes a prepared printable-ASCII input. Decoding and input generation are excluded.",
			Baseline:    "Base64",
			Operation:   "encoding",
			Note:        "The formats have different output expansion. Allocated bytes include buffer capacity and are not a measurement of encoded length.",
			Cases: []caseDef{
				{Label: "2 KiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_2KB"},
					{"ASCII85", "encoding", "BenchmarkBase85_2KB"},
				}},
				{Label: "4 KiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_4KB"},
					{"ASCII85", "encoding", "BenchmarkBase85_4KB"},
				}},
				{Label: "8 KiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_8KB"},
					{"ASCII85", "encoding", "BenchmarkBase85_8KB"},
				}},
				{Label: "16 KiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_16KB"},
					{"ASCII85", "encoding", "BenchmarkBase85_16KB"},
				}},
				{Label: "32 KiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_32KB"},
					{"ASCII85", "encoding", "BenchmarkBase85_32KB"},
				}},
				{Label: "64 KiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_64KB"},
					{"ASCII85", "encoding", "BenchmarkBase85_64KB"},
				}},
				{Label: "128 KiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_128KB"},
					{"ASCII85", "encoding", "BenchmarkBase85_128KB"},
				}},
				{Label: "256 KiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_256KB"},
					{"ASCII85", "encoding", "BenchmarkBase85_256KB"},
				}},
				{Label: "512 KiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_512KB"},
					{"ASCII85", "encoding", "BenchmarkBase85_512KB"},
				}},
				{Label: "1 MiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_1MB"},
					{"ASCII85", "encoding", "BenchmarkBase85_1MB"},
				}},
				{Label: "10 MiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_10MB"},
					{"ASCII85", "encoding", "BenchmarkBase85_10MB"},
				}},
				{Label: "100 MiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_100MB"},
					{"ASCII85", "encoding", "BenchmarkBase85_100MB"},
				}},
				{Label: "1 GiB", Rows: []rowDef{
					{"Base64", "encoding", "BenchmarkBase64_1GB"},
					{"ASCII85", "encoding", "BenchmarkBase85_1GB"},
				}},
			},
		},
	},
	"iteration": {
		{
			Title:       "Read all 100,000 rows",
			Description: "One operation creates an iterator, consumes every row, and closes it.",
			Baseline:    "Direct slice",
			Operation:   "scan",
			Note:        "The pull adapter uses iter.Pull; the direct-channel reference uses a producer goroutine and explicit channel handoff.",
			Rows: []rowDef{
				{"Direct slice", "iteration", "BenchmarkBaseline_SliceDirect_Full_100k"},
				{"Slice adapter", "iteration", "BenchmarkPullSlice_Full_100k"},
				{"Pull adapter", "iteration", "BenchmarkPullChan_Full_100k"},
				{"Direct channel", "iteration", "BenchmarkBaseline_ChanDirect_Full_100k"},
			},
		},
		{
			Title:       "Stop after 10 of 100,000 rows",
			Description: "Create an iterator, read ten rows, and close it.",
			Baseline:    "Slice adapter",
			Operation:   "scan",
			Note:        "The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.",
			Rows: []rowDef{
				{"Slice adapter", "iteration", "BenchmarkPullSlice_EarlyClose10_100k"},
				{"Pull adapter", "iteration", "BenchmarkPullChan_EarlyClose10_100k"},
			},
		},
		{
			Title:       "First row from 1 million rows",
			Description: "Create an iterator, read one row, and close it. This is the complete first-row lifecycle, including close.",
			Baseline:    "Slice adapter",
			Operation:   "scan",
			Note:        "The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.",
			Rows: []rowDef{
				{"Slice adapter", "iteration", "BenchmarkPullSlice_TTFR_1M"},
				{"Pull adapter", "iteration", "BenchmarkPullChan_TTFR_1M"},
			},
		},
		{
			Title:       "Consume 100,000 rows with consumer work",
			Description: "Create, drain, and close an iterator while running the workload's synthetic consumer loop for each row.",
			Baseline:    "Slice adapter",
			Operation:   "scan",
			Note:        "The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.",
			Rows: []rowDef{
				{"Slice adapter", "iteration", "BenchmarkPullSlice_SlowConsumerWork50_100k"},
				{"Pull adapter", "iteration", "BenchmarkPullChan_SlowConsumerWork50_100k"},
			},
		},
		{
			Title:       "Consume 100,000 rows with producer work",
			Description: "Create, drain, and close an iterator whose source runs the workload's synthetic producer loop for each row.",
			Baseline:    "Slice adapter",
			Operation:   "scan",
			Note:        "The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.",
			Rows: []rowDef{
				{"Slice adapter", "iteration", "BenchmarkPullSlice_SlowProducerWork50_100k"},
				{"Pull adapter", "iteration", "BenchmarkPullChan_SlowProducerWork50_100k"},
			},
		},
		{
			Title:       "Read 100,000 rows across ten ranges",
			Description: "Create an iterator over ten consecutive ranges, consume all rows, and close it.",
			Baseline:    "Slice adapter",
			Operation:   "scan",
			Note:        "The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.",
			Rows: []rowDef{
				{"Slice adapter", "iteration", "BenchmarkPullSlice_MultiRange_100k_10parts"},
				{"Pull adapter", "iteration", "BenchmarkPullChan_MultiRange_100k_10parts"},
			},
		},
		{
			Title:       "Open and close an empty range",
			Description: "Create an iterator for an empty source, request its end, and close it.",
			Baseline:    "Slice adapter",
			Operation:   "scan",
			Note:        "The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.",
			Rows: []rowDef{
				{"Slice adapter", "iteration", "BenchmarkPullSlice_EmptyRange"},
				{"Pull adapter", "iteration", "BenchmarkPullChan_EmptyRange"},
			},
		},
		{
			Title:       "Parallel scans of 10,000 rows",
			Description: "Each parallel operation creates an iterator, reads 10,000 rows, and closes it; parallelism is four times GOMAXPROCS.",
			Baseline:    "Slice adapter",
			Operation:   "scan",
			Note:        "The slice adapter preloads the requested range; the pull adapter produces rows on demand. Allocation totals do not measure peak retained memory.",
			Rows: []rowDef{
				{"Slice adapter", "iteration", "BenchmarkPullSlice_Parallel_10k"},
				{"Pull adapter", "iteration", "BenchmarkPullChan_Parallel_10k"},
			},
		},
	},
	"log": {
		{
			Title:       "Simple INFO message",
			Description: "Encode one INFO record with a timestamp and no attributes. All output goes to io.Discard.",
			Baseline:    "slog",
			Operation:   "operation",
			Note:        "Each logger keeps its native JSON field names and timestamp format. Source locations are disabled.",
			Rows: []rowDef{
				{"slog", "log", "BenchmarkSlog_SimpleMessage"},
				{"slogx", "log", "BenchmarkSlogx_SimpleMessage"},
				{"zerolog", "log", "BenchmarkZerolog_SimpleMessage"},
				{"zap", "log", "BenchmarkZap_SimpleMessage"},
				{"logrus", "log", "BenchmarkLogrus_SimpleMessage"},
			},
		},
		{
			Title:       "INFO message with three fields",
			Description: "Encode one INFO record with a timestamp and three structured attributes. All output goes to io.Discard.",
			Baseline:    "slog",
			Operation:   "operation",
			Note:        "Each logger keeps its native JSON field names and timestamp format. Source locations are disabled.",
			Rows: []rowDef{
				{"slog", "log", "BenchmarkSlog_WithFields"},
				{"slogx", "log", "BenchmarkSlogx_WithFields"},
				{"zerolog", "log", "BenchmarkZerolog_WithFields"},
				{"zap", "log", "BenchmarkZap_WithFields"},
				{"logrus", "log", "BenchmarkLogrus_WithFields"},
			},
		},
		{
			Title:       "Plain error",
			Description: "Encode one ERROR record carrying an existing plain error. All output goes to io.Discard.",
			Baseline:    "slog",
			Operation:   "operation",
			Note:        "Each logger keeps its native JSON field names and timestamp format. Source locations are disabled.",
			Rows: []rowDef{
				{"slog", "log", "BenchmarkSlog_Error"},
				{"slogx", "log", "BenchmarkSlogx_Error"},
				{"zerolog", "log", "BenchmarkZerolog_Error"},
				{"zap", "log", "BenchmarkZap_Error"},
				{"logrus", "log", "BenchmarkLogrus_Error"},
			},
		},
		{
			Title:       "Error with three fields",
			Description: "Encode one ERROR record carrying an existing plain error and three attributes. All output goes to io.Discard.",
			Baseline:    "slog",
			Operation:   "operation",
			Note:        "Each logger keeps its native JSON field names and timestamp format. Source locations are disabled.",
			Rows: []rowDef{
				{"slog", "log", "BenchmarkSlog_ErrorWithFields"},
				{"slogx", "log", "BenchmarkSlogx_ErrorWithFields"},
				{"zerolog", "log", "BenchmarkZerolog_ErrorWithFields"},
				{"zap", "log", "BenchmarkZap_ErrorWithFields"},
				{"logrus", "log", "BenchmarkLogrus_ErrorWithFields"},
			},
		},
		{
			Title:       "Three correlated error records",
			Description: "Encode three ERROR records sharing a prepared trace ID; the first record carries an existing plain error. Identifier generation is excluded. All output goes to io.Discard.",
			Baseline:    "slog",
			Operation:   "operation",
			Note:        "Each logger keeps its native JSON field names and timestamp format. Source locations are disabled.",
			Rows: []rowDef{
				{"slog", "log", "BenchmarkSlog_ErrorChain"},
				{"slogx", "log", "BenchmarkSlogx_ErrorChain"},
				{"zerolog", "log", "BenchmarkZerolog_ErrorChain"},
				{"zap", "log", "BenchmarkZap_ErrorChain"},
				{"logrus", "log", "BenchmarkLogrus_ErrorChain"},
			},
		},
		{
			Title:       "slogx context enrichment",
			Description: "Encode one INFO record with a call-site attribute using a prepared context. Context construction is excluded.",
			Baseline:    "No context fields",
			Operation:   "record",
			Note:        "Ratios describe the additional enrichment cost; these cases intentionally emit different fields.",
			Rows: []rowDef{
				{"No context fields", "log", "BenchmarkSlogxFeatures_Context/NoAttrs"},
				{"3 context fields", "log", "BenchmarkSlogxFeatures_Context/Attrs3"},
				{"10 context fields", "log", "BenchmarkSlogxFeatures_Context/Attrs10"},
				{"Trace and span identity", "log", "BenchmarkSlogxFeatures_Context/TraceOnly"},
			},
		},
		{
			Title:       "slogx recording spans",
			Description: "Start and end a span with two attributes, including identifier generation, context allocation, and JSON emission.",
			Operation:   "span",
			Note:        "Root and child spans have different parent-context work. No equivalent alternative implementation was measured.",
			Rows: []rowDef{
				{"Root span", "log", "BenchmarkSlogxFeatures_Span/RootRecording"},
				{"Child of a prepared parent", "log", "BenchmarkSlogxFeatures_Span/ChildRecording"},
			},
		},
		{
			Title:       "slogx non-recording spans",
			Description: "Start and end a span without emitting its JSON record. Identifier generation, context allocation, and lifecycle work remain included.",
			Operation:   "span",
			Note:        "These are distinct suppression and parent-context paths, kept separate from recording spans.",
			Rows: []rowDef{
				{"Root, level disabled", "log", "BenchmarkSlogxFeatures_Span/RootLevelDisabled"},
				{"Child, level disabled", "log", "BenchmarkSlogxFeatures_Span/ChildLevelDisabled"},
				{"Child, parent unsampled", "log", "BenchmarkSlogxFeatures_Span/ChildUnsampled"},
			},
		},
		{
			Title:       "slogx structured error lifecycle",
			Description: "Construct three attributed error wrappers, or encode a prepared three-layer error in one ERROR record.",
			Operation:   "operation",
			Note:        "Construction excludes logging; logging excludes construction. Neither row is a baseline for the other.",
			Rows: []rowDef{
				{"Construct three wrappers", "log", "BenchmarkSlogxFeatures_StructuredError/WrapThreeLayers"},
				{"Log a prepared wrapped error", "log", "BenchmarkSlogxFeatures_StructuredError/LogThreeLayers"},
			},
		},
	},
	"meta": {
		{
			Title:       "Sort random integers",
			Description: "One operation allocates and copies the input slice, then sorts it in ascending order. Each implementation receives a separately generated input of the same size and distribution.",
			Baseline:    "Go sort.Ints",
			Operation:   "sort",
			Note:        "Libraries retain their default worker policies; Sorty uses MaxGor=3. GOMAXPROCS does not mean every implementation uses that many workers.",
			Cases: []caseDef{
				{Label: "10,000 integers", Rows: []rowDef{
					{"Go sort.Ints", "meta", "BenchmarkSortInts/Asc_Int_10000"},
					{"Sorty", "meta", "BenchmarkSortyInts/Asc_Int_10000"},
					{"Pargo", "meta", "BenchmarkPargoInts/Asc_Int_10000"},
					{"Parsort", "meta", "BenchmarkParsortInts/Asc_Int_10000"},
				}},
				{Label: "100,000 integers", Rows: []rowDef{
					{"Go sort.Ints", "meta", "BenchmarkSortInts/Asc_Int_100000"},
					{"Sorty", "meta", "BenchmarkSortyInts/Asc_Int_100000"},
					{"Pargo", "meta", "BenchmarkPargoInts/Asc_Int_100000"},
					{"Parsort", "meta", "BenchmarkParsortInts/Asc_Int_100000"},
				}},
				{Label: "1,000,000 integers", Rows: []rowDef{
					{"Go sort.Ints", "meta", "BenchmarkSortInts/Asc_Int_1000000"},
					{"Sorty", "meta", "BenchmarkSortyInts/Asc_Int_1000000"},
					{"Pargo", "meta", "BenchmarkPargoInts/Asc_Int_1000000"},
					{"Parsort", "meta", "BenchmarkParsortInts/Asc_Int_1000000"},
				}},
				{Label: "10,000,000 integers", Rows: []rowDef{
					{"Go sort.Ints", "meta", "BenchmarkSortInts/Asc_Int_10000000"},
					{"Sorty", "meta", "BenchmarkSortyInts/Asc_Int_10000000"},
					{"Pargo", "meta", "BenchmarkPargoInts/Asc_Int_10000000"},
					{"Parsort", "meta", "BenchmarkParsortInts/Asc_Int_10000000"},
				}},
			},
		},
		{
			Title:       "Set one string field",
			Description: "Set one string field through its direct method or a prepared reflect.Value.",
			Baseline:    "Direct method",
			Operation:   "access",
			Note:        "Field lookup and reflect.Value preparation happen before timing; this measures field access only.",
			Rows: []rowDef{
				{"Direct method", "meta", "BenchmarkMethodAccess/SetName"},
				{"Reflection", "meta", "BenchmarkReflectionAccess/SetName"},
			},
		},
		{
			Title:       "Read one string field",
			Description: "Read one string field through its direct method or a prepared reflect.Value.",
			Baseline:    "Direct method",
			Operation:   "access",
			Note:        "Field lookup and reflect.Value preparation happen before timing; this measures field access only.",
			Rows: []rowDef{
				{"Direct method", "meta", "BenchmarkMethodAccess/GetName"},
				{"Reflection", "meta", "BenchmarkReflectionAccess/GetName"},
			},
		},
		{
			Title:       "Synchronous shared-counter increments",
			Description: "Parallel workers increment one exact shared count and wait for each update to complete.",
			Baseline:    "Atomic",
			Operation:   "increment",
			Note:        "Parallel time per operation is elapsed time divided by all workers' operations, not individual-call latency. CAS backoff includes scheduler yields.",
			Rows: []rowDef{
				{"Atomic", "meta", "BenchmarkCounters_IncOnly/Atomic"},
				{"Mutex", "meta", "BenchmarkCounters_IncOnly/Mutex"},
				{"RWMutex", "meta", "BenchmarkCounters_IncOnly/RWMutex"},
				{"CAS spin", "meta", "BenchmarkCounters_IncOnly/CAS_Spin"},
				{"CAS with backoff", "meta", "BenchmarkCounters_IncOnly/CAS_Backoff"},
				{"Flat combining", "meta", "BenchmarkReduction_IncOnly/FlatCombining"},
				{"RCU copy-on-write", "meta", "BenchmarkReduction_IncOnly/RCU"},
			},
		},
		{
			Title:       "Shared-counter reads",
			Description: "Parallel workers read a shared counter without concurrent writers.",
			Baseline:    "Atomic",
			Operation:   "read",
			Note:        "Parallel time per operation is elapsed time divided by all workers' operations, not individual-call latency. CAS backoff includes scheduler yields.",
			Rows: []rowDef{
				{"Atomic", "meta", "BenchmarkCounters_GetOnly/Atomic"},
				{"Mutex", "meta", "BenchmarkCounters_GetOnly/Mutex"},
				{"RWMutex", "meta", "BenchmarkCounters_GetOnly/RWMutex"},
				{"CAS spin", "meta", "BenchmarkCounters_GetOnly/CAS_Spin"},
				{"CAS with backoff", "meta", "BenchmarkCounters_GetOnly/CAS_Backoff"},
				{"RCU pointer read", "meta", "BenchmarkReduction_GetOnly/RCU"},
			},
		},
		{
			Title:       "Shared counter: 90% updates, 10% reads",
			Description: "Each worker repeats nine completed increments and one read. One reported operation is one increment or read, not a ten-operation batch.",
			Baseline:    "Atomic",
			Operation:   "operation",
			Note:        "Parallel time per operation is elapsed time divided by all workers' operations, not individual-call latency. CAS backoff includes scheduler yields.",
			Rows: []rowDef{
				{"Atomic", "meta", "BenchmarkCounters_Mixed90_10/Atomic"},
				{"Mutex", "meta", "BenchmarkCounters_Mixed90_10/Mutex"},
				{"RWMutex", "meta", "BenchmarkCounters_Mixed90_10/RWMutex"},
				{"CAS spin", "meta", "BenchmarkCounters_Mixed90_10/CAS_Spin"},
				{"CAS with backoff", "meta", "BenchmarkCounters_Mixed90_10/CAS_Backoff"},
			},
		},
		{
			Title:       "Counter increment submission by channel capacity",
			Description: "Parallel workers submit increment requests to one owner goroutine; return does not confirm that an increment has been applied.",
			Baseline:    "Unbuffered",
			Operation:   "operation",
			Note:        "These owner-goroutine operations have different completion boundaries from atomic and locked counter calls.",
			Rows: []rowDef{
				{"Unbuffered", "meta", "BenchmarkCounters_IncOnly/Channel_buf0"},
				{"Capacity 1", "meta", "BenchmarkCounters_IncOnly/Channel_buf1"},
				{"Capacity 256", "meta", "BenchmarkCounters_IncOnly/Channel_buf256"},
				{"Capacity 4,096", "meta", "BenchmarkCounters_IncOnly/Channel_buf4096"},
			},
		},
		{
			Title:       "Counter request/reply reads by channel capacity",
			Description: "Parallel workers send read requests to an owner goroutine and wait for their replies.",
			Baseline:    "Unbuffered",
			Operation:   "operation",
			Note:        "These owner-goroutine operations have different completion boundaries from atomic and locked counter calls.",
			Rows: []rowDef{
				{"Unbuffered", "meta", "BenchmarkCounters_GetOnly/Channel_buf0"},
				{"Capacity 1", "meta", "BenchmarkCounters_GetOnly/Channel_buf1"},
				{"Capacity 256", "meta", "BenchmarkCounters_GetOnly/Channel_buf256"},
				{"Capacity 4,096", "meta", "BenchmarkCounters_GetOnly/Channel_buf4096"},
			},
		},
		{
			Title:       "Channel counter: 90% submissions, 10% reads",
			Description: "Each worker repeats nine increment submissions and one request/reply read. One operation is one submission or read.",
			Baseline:    "Unbuffered",
			Operation:   "operation",
			Note:        "These owner-goroutine operations have different completion boundaries from atomic and locked counter calls.",
			Rows: []rowDef{
				{"Unbuffered", "meta", "BenchmarkCounters_Mixed90_10/Channel_buf0"},
				{"Capacity 1", "meta", "BenchmarkCounters_Mixed90_10/Channel_buf1"},
				{"Capacity 256", "meta", "BenchmarkCounters_Mixed90_10/Channel_buf256"},
				{"Capacity 4,096", "meta", "BenchmarkCounters_Mixed90_10/Channel_buf4096"},
			},
		},
		{
			Title:       "Single-producer, single-consumer transfer",
			Description: "One operation transfers a value between one producer and one consumer. Both implementations retry full/empty operations with scheduler yields.",
			Baseline:    "Channel",
			Operation:   "transfer",
			Note:        "A ring buffer reserves one slot, so its usable capacity is one less than the listed channel capacity.",
			Cases: []caseDef{
				{Label: "256 slots", Rows: []rowDef{
					{"Channel", "meta", "BenchmarkSPSC/Chan_buf256"},
					{"Ring buffer", "meta", "BenchmarkSPSC/RingBuffer_size256"},
				}},
				{Label: "4,096 slots", Rows: []rowDef{
					{"Channel", "meta", "BenchmarkSPSC/Chan_buf4096"},
					{"Ring buffer", "meta", "BenchmarkSPSC/RingBuffer_size4096"},
				}},
				{Label: "65,536 slots", Rows: []rowDef{
					{"Channel", "meta", "BenchmarkSPSC/Chan_buf65536"},
					{"Ring buffer", "meta", "BenchmarkSPSC/RingBuffer_size65536"},
				}},
			},
		},
		{
			Title:       "Concurrent map stores",
			Description: "Parallel workers store values across 1,024 shared string keys; timing includes first insertion and subsequent replacement.",
			Baseline:    "Mutex map",
			Operation:   "access",
			Note:        "Results describe this shared-key workload; time per access is an aggregate-throughput measure.",
			Rows: []rowDef{
				{"Mutex map", "meta", "BenchmarkMap_StoreOnly/MutexMap"},
				{"sync.Map", "meta", "BenchmarkMap_StoreOnly/SyncMap"},
			},
		},
		{
			Title:       "Concurrent map loads",
			Description: "Parallel workers load values across 1,024 shared string keys populated before timing.",
			Baseline:    "Mutex map",
			Operation:   "access",
			Note:        "Results describe this shared-key workload; time per access is an aggregate-throughput measure.",
			Rows: []rowDef{
				{"Mutex map", "meta", "BenchmarkMap_LoadOnly/MutexMap"},
				{"sync.Map", "meta", "BenchmarkMap_LoadOnly/SyncMap"},
			},
		},
		{
			Title:       "Sharded counter updates",
			Description: "Parallel workers increment their selected padded atomic shard or slot.",
			Baseline:    "8 shards",
			Operation:   "increment",
			Note:        "Workers are assigned by ID, not pinned to physical CPUs. Reading the sum is separate work and is not an atomic snapshot during concurrent writes.",
			Rows: []rowDef{
				{"8 shards", "meta", "BenchmarkReduction_IncOnly/Sharded_8"},
				{"One shard per logical CPU", "meta", "BenchmarkReduction_IncOnly/Sharded_NumCPU"},
				{"GOMAXPROCS slots", "meta", "BenchmarkReduction_IncOnly/PerCPU"},
			},
		},
		{
			Title:       "Read a sharded total",
			Description: "Parallel readers sum padded atomic slots without concurrent writers.",
			Baseline:    "8 shards",
			Operation:   "sum",
			Note:        "Larger shard arrays change the amount of read work. These measurements do not establish snapshot consistency with active writers.",
			Rows: []rowDef{
				{"8 shards", "meta", "BenchmarkReduction_GetOnly/Sharded_8"},
				{"GOMAXPROCS slots", "meta", "BenchmarkReduction_GetOnly/PerCPU"},
			},
		},
		{
			Title:       "Buffered local increments",
			Description: "Workers accumulate increments locally and publish them at the configured threshold, then flush their remainders before finishing.",
			Baseline:    "Flush every 16",
			Operation:   "increment",
			Note:        "A larger threshold changes publication frequency and visibility. Reads omit unflushed increments.",
			Rows: []rowDef{
				{"Flush every 16", "meta", "BenchmarkReduction_IncOnly/LocalBuffered_t16"},
				{"Flush every 64", "meta", "BenchmarkReduction_IncOnly/LocalBuffered_t64"},
				{"Flush every 256", "meta", "BenchmarkReduction_IncOnly/LocalBuffered_t256"},
			},
		},
		{
			Title:       "MPSC counter submission",
			Description: "Parallel producers submit increments through a buffered channel to one aggregator.",
			Baseline:    "Capacity 256",
			Operation:   "submission",
			Note:        "Submission does not wait for the aggregate total to be published. These are not completed counter-update timings.",
			Rows: []rowDef{
				{"Capacity 256", "meta", "BenchmarkReduction_IncOnly/MPSC_buf256"},
				{"Capacity 4,096", "meta", "BenchmarkReduction_IncOnly/MPSC_buf4096"},
			},
		},
		{
			Title:       "Approximate counter update costs",
			Description: "Workers sample updates with a local pseudo-random state; an accepted update adds the configured sampling rate.",
			Operation:   "attempt",
			Note:        "Sampling rates trade execution cost for estimation error. Estimator accuracy was not measured, so no relative-speed ranking is assigned.",
			Rows: []rowDef{
				{"Sample 1 in 10", "meta", "BenchmarkReduction_IncOnly/Approx_rate10"},
				{"Sample 1 in 100", "meta", "BenchmarkReduction_IncOnly/Approx_rate100"},
				{"Sample 1 in 1,000", "meta", "BenchmarkReduction_IncOnly/Approx_rate1000"},
			},
		},
		{
			Title:       "Read an approximate counter",
			Description: "Parallel readers load the current estimate from a counter configured to sample one update in 100; there are no concurrent writers.",
			Operation:   "read",
			Note:        "This reads an estimate; no accuracy-equivalent alternative was measured.",
			Rows: []rowDef{
				{"Sample rate 100", "meta", "BenchmarkReduction_GetOnly/Approx_rate100"},
			},
		},
		{
			Title:       "Four-producer queue transfer",
			Description: "Four producers publish messages and one consumer receives all of them through a Disruptor-style ring. One operation is one complete transfer.",
			Baseline:    "256 slots",
			Operation:   "transfer",
			Note:        "Timing includes goroutine startup, retries, and completion waiting. This is separate from channel submission-only measurements.",
			Rows: []rowDef{
				{"256 slots", "meta", "BenchmarkDisruptorMPSC/size256"},
				{"4,096 slots", "meta", "BenchmarkDisruptorMPSC/size4096"},
				{"65,536 slots", "meta", "BenchmarkDisruptorMPSC/size65536"},
			},
		},
		{
			Title:       "Semaphore acquire and release",
			Description: "Parallel workers acquire and release a channel-based semaphore without any protected work between the calls.",
			Baseline:    "Capacity 1",
			Operation:   "acquire/release",
			Note:        "Capacity changes the permitted concurrency. This isolates semaphore overhead, not useful protected-work throughput.",
			Rows: []rowDef{
				{"Capacity 1", "meta", "BenchmarkSemaphore/max1"},
				{"Capacity 4", "meta", "BenchmarkSemaphore/max4"},
				{"Capacity 16", "meta", "BenchmarkSemaphore/max16"},
				{"Capacity 64", "meta", "BenchmarkSemaphore/max64"},
			},
		},
		{
			Title:       "Start goroutines and wait",
			Description: "One operation creates the listed number of goroutines with empty work and waits for all of them to finish.",
			Operation:   "batch",
			Note:        "Rows perform different amounts of startup work; the time is per complete batch, not per goroutine.",
			Rows: []rowDef{
				{"1 goroutine", "meta", "BenchmarkSpawnAndWait/workers1"},
				{"4 goroutines", "meta", "BenchmarkSpawnAndWait/workers4"},
				{"16 goroutines", "meta", "BenchmarkSpawnAndWait/workers16"},
				{"64 goroutines", "meta", "BenchmarkSpawnAndWait/workers64"},
			},
		},
		{
			Title:       "Pool get and put",
			Description: "Get a reusable buffer from sync.Pool and return it, serially or across parallel workers.",
			Operation:   "get/put",
			Note:        "The serial case measures a sequential operation; the parallel case reports elapsed time divided by total worker operations. Their ratio is not per-call speedup.",
			Rows: []rowDef{
				{"Serial", "meta", "BenchmarkPool/Serial"},
				{"Parallel", "meta", "BenchmarkPool/Parallel"},
			},
		},
		{
			Title:       "Once initialization paths",
			Description: "Initialize and read a fresh instance, reuse one shared instance across workers, or reuse an instance initialized before timing.",
			Operation:   "operation",
			Note:        "The shared instance initializes only once across the benchmark. These initialization states and operation boundaries are different.",
			Rows: []rowDef{
				{"Fresh instance, serial", "meta", "BenchmarkOnce/Serial"},
				{"Shared instance, parallel", "meta", "BenchmarkOnce/Contended"},
				{"Preinitialized, parallel", "meta", "BenchmarkOnce/FastPath_AfterInit"},
			},
		},
	},
	"protocol": {
		{
			Title:       "HTTP payload scaling",
			Description: "POST a prepared body to a loopback server and wait for HTTP 200 after the handler reads the body.",
			Operation:   "message",
			Note:        "Completion includes server acknowledgement. These costs are not normalized against TCP writes, which have no application acknowledgement. Go allocations can include the server goroutines in the same process.",
			Cases: []caseDef{
				{Label: "9 bytes", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage"},
				}},
				{Label: "2 KiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage2KB"},
				}},
				{Label: "4 KiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage4KB"},
				}},
				{Label: "8 KiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage8KB"},
				}},
				{Label: "16 KiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage16KB"},
				}},
				{Label: "32 KiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage32KB"},
				}},
				{Label: "64 KiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage64KB"},
				}},
				{Label: "128 KiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage128KB"},
				}},
				{Label: "256 KiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage256KB"},
				}},
				{Label: "512 KiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage512KB"},
				}},
				{Label: "1 MiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage1MB"},
				}},
				{Label: "10 MiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage10MB"},
				}},
				{Label: "100 MiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage100MB"},
				}},
				{Label: "1 GiB", Rows: []rowDef{
					{"HTTP", "protocol/http", "BenchmarkHTTPClientSendSingleMessage1GB"},
				}},
			},
		},
		{
			Title:       "TCP payload scaling",
			Description: "Write and flush a prepared payload to a loopback socket while a server drains the stream.",
			Operation:   "message",
			Note:        "Completion does not confirm server consumption of the full message. These costs are not normalized against HTTP request/response completion. Go allocations can include the server goroutines in the same process.",
			Cases: []caseDef{
				{Label: "9 bytes", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage"},
				}},
				{Label: "2 KiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage2KB"},
				}},
				{Label: "4 KiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage4KB"},
				}},
				{Label: "8 KiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage8KB"},
				}},
				{Label: "16 KiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage16KB"},
				}},
				{Label: "32 KiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage32KB"},
				}},
				{Label: "64 KiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage64KB"},
				}},
				{Label: "128 KiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage128KB"},
				}},
				{Label: "256 KiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage256KB"},
				}},
				{Label: "512 KiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage512KB"},
				}},
				{Label: "1 MiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage1MB"},
				}},
				{Label: "10 MiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage10MB"},
				}},
				{Label: "100 MiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage100MB"},
				}},
				{Label: "1 GiB", Rows: []rowDef{
					{"TCP", "protocol/tcp", "BenchmarkTCPClientSendSingleMessage1GB"},
				}},
			},
		},
	},
	"serializer": {
		{
			Title:       "Small records: fresh round trips",
			Description: "One operation serializes and deserializes the listed number of generated Person records. Input generation and Protobuf type conversion are excluded.",
			Baseline:    "JSON",
			Operation:   "round trip",
			Note:        "JSON and Gob process a slice; Protobuf processes each message separately. The ratios compare these implemented workload choices, including their different framing.",
			Cases: []caseDef{
				{Label: "1 record", Rows: []rowDef{
					{"JSON", "serializer", "BenchmarkEncodeDecodeJson1Small"},
					{"Gob, fresh stream", "serializer", "BenchmarkEncodeDecodeBytes1Small"},
					{"Protobuf", "serializer", "BenchmarkEncodeDecodeProto1Small"},
				}},
				{Label: "100 records", Rows: []rowDef{
					{"JSON", "serializer", "BenchmarkEncodeDecodeJson100Small"},
					{"Gob, fresh stream", "serializer", "BenchmarkEncodeDecodeBytes100Small"},
					{"Protobuf", "serializer", "BenchmarkEncodeDecodeProto100Small"},
				}},
				{Label: "10,000 records", Rows: []rowDef{
					{"JSON", "serializer", "BenchmarkEncodeDecodeJson10000Small"},
					{"Gob, fresh stream", "serializer", "BenchmarkEncodeDecodeBytes10000Small"},
					{"Protobuf", "serializer", "BenchmarkEncodeDecodeProto10000Small"},
				}},
				{Label: "1,000,000 records", Rows: []rowDef{
					{"JSON", "serializer", "BenchmarkEncodeDecodeJson1000000Small"},
					{"Gob, fresh stream", "serializer", "BenchmarkEncodeDecodeBytes1000000Small"},
					{"Protobuf", "serializer", "BenchmarkEncodeDecodeProto1000000Small"},
				}},
			},
		},
		{
			Title:       "Small records: warmed Gob stream",
			Description: "One operation encodes and decodes the listed number of Person records using a reused buffer, encoder, and decoder after stream metadata has been warmed.",
			Operation:   "round trip",
			Note:        "Stream metadata is amortized across operations. These bytes are not independently decodable standalone messages, so this panel is separate from fresh round trips.",
			Cases: []caseDef{
				{Label: "1 record", Rows: []rowDef{
					{"Gob, warmed stream", "serializer", "BenchmarkEncodeDecodeBytesOpt1Small"},
				}},
				{Label: "100 records", Rows: []rowDef{
					{"Gob, warmed stream", "serializer", "BenchmarkEncodeDecodeBytesOpt100Small"},
				}},
				{Label: "10,000 records", Rows: []rowDef{
					{"Gob, warmed stream", "serializer", "BenchmarkEncodeDecodeBytesOpt10000Small"},
				}},
				{Label: "1,000,000 records", Rows: []rowDef{
					{"Gob, warmed stream", "serializer", "BenchmarkEncodeDecodeBytesOpt1000000Small"},
				}},
			},
		},
		{
			Title:       "Large records: fresh round trips",
			Description: "One operation serializes and deserializes the listed number of generated 100-field Unreal records. Input generation and Protobuf type conversion are excluded.",
			Baseline:    "JSON",
			Operation:   "round trip",
			Note:        "JSON and Gob process a slice; Protobuf processes each message separately. The ratios compare these implemented workload choices, including their different framing.",
			Cases: []caseDef{
				{Label: "1 record", Rows: []rowDef{
					{"JSON", "serializer", "BenchmarkEncodeDecodeJson1Unreal"},
					{"Gob, fresh stream", "serializer", "BenchmarkEncodeDecodeBytes1Unreal"},
					{"Protobuf", "serializer", "BenchmarkEncodeDecodeProto1Unreal"},
				}},
				{Label: "10 records", Rows: []rowDef{
					{"JSON", "serializer", "BenchmarkEncodeDecodeJson10Unreals"},
					{"Gob, fresh stream", "serializer", "BenchmarkEncodeDecodeBytes10Unreals"},
					{"Protobuf", "serializer", "BenchmarkEncodeDecodeProto10Unreals"},
				}},
				{Label: "100 records", Rows: []rowDef{
					{"JSON", "serializer", "BenchmarkEncodeDecodeJson100Unreals"},
					{"Gob, fresh stream", "serializer", "BenchmarkEncodeDecodeBytes100Unreals"},
					{"Protobuf", "serializer", "BenchmarkEncodeDecodeProto100Unreals"},
				}},
				{Label: "1,000 records", Rows: []rowDef{
					{"JSON", "serializer", "BenchmarkEncodeDecodeJson1000Unreals"},
					{"Gob, fresh stream", "serializer", "BenchmarkEncodeDecodeBytes1000Unreals"},
					{"Protobuf", "serializer", "BenchmarkEncodeDecodeProto1000Unreals"},
				}},
			},
		},
		{
			Title:       "Large records: warmed Gob stream",
			Description: "One operation encodes and decodes the listed number of 100-field Unreal records using a reused buffer, encoder, and decoder after stream metadata has been warmed.",
			Operation:   "round trip",
			Note:        "Stream metadata is amortized across operations. These bytes are not independently decodable standalone messages, so this panel is separate from fresh round trips.",
			Cases: []caseDef{
				{Label: "1 record", Rows: []rowDef{
					{"Gob, warmed stream", "serializer", "BenchmarkEncodeDecodeBytesOpt1Unreal"},
				}},
				{Label: "10 records", Rows: []rowDef{
					{"Gob, warmed stream", "serializer", "BenchmarkEncodeDecodeBytesOpt10Unreals"},
				}},
				{Label: "100 records", Rows: []rowDef{
					{"Gob, warmed stream", "serializer", "BenchmarkEncodeDecodeBytesOpt100Unreals"},
				}},
				{Label: "1,000 records", Rows: []rowDef{
					{"Gob, warmed stream", "serializer", "BenchmarkEncodeDecodeBytesOpt1000Unreals"},
				}},
			},
		},
	},
	"validator": {
		{
			Title:       "One valid email",
			Description: "Validate one valid email.",
			Baseline:    "Playground",
			Operation:   "validation",
			Note:        "Playground reuses its validator and caches; Ward constructs validators per operation; Ozzo uses its stateless API. Native email and URL checks can differ beyond these fixtures.",
			Rows: []rowDef{
				{"Playground", "validator", "BenchmarkPlaygroundSingleFieldValid"},
				{"Ward", "validator", "BenchmarkWardSingleFieldValid"},
				{"Ozzo", "validator", "BenchmarkOzzoSingleFieldValid"},
			},
		},
		{
			Title:       "One invalid email",
			Description: "Validate one invalid email and retain the native failure result.",
			Baseline:    "Playground",
			Operation:   "validation",
			Note:        "Playground reuses its validator and caches; Ward constructs validators per operation; Ozzo uses its stateless API. Native email and URL checks can differ beyond these fixtures.",
			Rows: []rowDef{
				{"Playground", "validator", "BenchmarkPlaygroundSingleFieldInvalid"},
				{"Ward", "validator", "BenchmarkWardSingleFieldInvalid"},
				{"Ozzo", "validator", "BenchmarkOzzoSingleFieldInvalid"},
			},
		},
		{
			Title:       "Five valid fields",
			Description: "Validate five valid string fields with the declared matching fixture rules.",
			Baseline:    "Playground",
			Operation:   "validation",
			Note:        "Playground reuses its validator and caches; Ward constructs validators per operation; Ozzo uses its stateless API. Native email and URL checks can differ beyond these fixtures.",
			Rows: []rowDef{
				{"Playground", "validator", "BenchmarkPlaygroundMultiFieldAllValid"},
				{"Ward", "validator", "BenchmarkWardMultiFieldAllValid"},
				{"Ozzo", "validator", "BenchmarkOzzoMultiFieldAllValid"},
			},
		},
		{
			Title:       "Five invalid fields",
			Description: "Validate the invalid fixture through all five fields; every field fails.",
			Baseline:    "Playground",
			Operation:   "validation",
			Note:        "Playground reuses its validator and caches; Ward constructs validators per operation; Ozzo uses its stateless API. Native email and URL checks can differ beyond these fixtures.",
			Rows: []rowDef{
				{"Playground", "validator", "BenchmarkPlaygroundMultiFieldSomeInvalid"},
				{"Ward", "validator", "BenchmarkWardMultiFieldSomeInvalid"},
				{"Ozzo", "validator", "BenchmarkOzzoMultiFieldSomeInvalid"},
			},
		},
		{
			Title:       "One valid email: validate and inspect",
			Description: "Validate one valid email. Check the success result.",
			Baseline:    "Playground",
			Operation:   "validation",
			Note:        "All three cases validate and inspect success, with their declared initialization/reuse policies.",
			Rows: []rowDef{
				{"Playground", "validator", "BenchmarkPlaygroundFullCycleSingleFieldValid"},
				{"Ward", "validator", "BenchmarkWardFullCycleSingleFieldValid"},
				{"Ozzo", "validator", "BenchmarkOzzoFullCycleSingleFieldValid"},
			},
		},
		{
			Title:       "One invalid email: validate and inspect",
			Description: "Validate one invalid email and retain the native failure result. Inspect the library's native failure representation.",
			Operation:   "validation",
			Note:        "Ozzo formats error text; Playground and Ward inspect field/rule identifiers. These consumer operations are not equivalent, so no speed ratio is assigned.",
			Rows: []rowDef{
				{"Playground", "validator", "BenchmarkPlaygroundFullCycleSingleFieldInvalid"},
				{"Ward", "validator", "BenchmarkWardFullCycleSingleFieldInvalid"},
				{"Ozzo", "validator", "BenchmarkOzzoFullCycleSingleFieldInvalid"},
			},
		},
		{
			Title:       "Five valid fields: validate and inspect",
			Description: "Validate five valid string fields with the declared matching fixture rules. Check the success result.",
			Baseline:    "Playground",
			Operation:   "validation",
			Note:        "All three cases validate and inspect success, with their declared initialization/reuse policies.",
			Rows: []rowDef{
				{"Playground", "validator", "BenchmarkPlaygroundFullCycleMultiFieldAllValid"},
				{"Ward", "validator", "BenchmarkWardFullCycleMultiFieldAllValid"},
				{"Ozzo", "validator", "BenchmarkOzzoFullCycleMultiFieldAllValid"},
			},
		},
		{
			Title:       "Five invalid fields: validate and inspect",
			Description: "Validate the invalid fixture through all five fields; every field fails. Inspect the library's native failure representation.",
			Operation:   "validation",
			Note:        "Ozzo formats error text; Playground and Ward inspect field/rule identifiers. These consumer operations are not equivalent, so no speed ratio is assigned.",
			Rows: []rowDef{
				{"Playground", "validator", "BenchmarkPlaygroundFullCycleMultiFieldSomeInvalid"},
				{"Ward", "validator", "BenchmarkWardFullCycleMultiFieldSomeInvalid"},
				{"Ozzo", "validator", "BenchmarkOzzoFullCycleMultiFieldSomeInvalid"},
			},
		},
		{
			Title:       "Ward: stop at the first failure",
			Description: "Validate the shared invalid fixture and stop after its first failing field, the email.",
			Operation:   "validation",
			Note:        "This checks less of the form than the five-field failure comparison. No equivalent stop-on-first-field alternative was measured.",
			Rows: []rowDef{
				{"Ward, stop on failure", "validator", "BenchmarkWardMultiFieldStopOnFail"},
			},
		},
	},
}
