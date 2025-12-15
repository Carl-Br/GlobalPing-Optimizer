# Globalping Optimizer
A Go-based tool that leverages the [Globalping.io API](https://globalping.io) to analyze HTTP request latencies from multiple geographic locations, helping you identify the optimal server location for your websites or REST APIs. 
Unlike traditional latency testing tools, it can identify optimal server placement even when the actual server location is hidden behind CDNs (like Cloudflare, Cloudfront, etc.).
## Demo
https://github.com/user-attachments/assets/ab22f6a4-f158-4fc2-87ef-88b32daac7ab
## Key Features
- **Multi-Location Testing**: Test from multiple geographic locations simultaneously
- **Flexible Configuration**: Support for various location formats including countries, cities, regions, ASNs, and cloud providers
- **Automated Measurements**: Run multiple measurement rounds with configurable intervals
- **Statistics**: Export comprehensive results to CSV format for further analysis
## Installation
### Prerequisites
- Go 1.25.3 or higher
- A Globalping.io [API token](https://dash.globalping.io/tokens)
### Build from Source
```bash
git clone https://github.com/Carl-Br/GlobalPing-Optimizer.git
cd globalping-optimizer
```
## Configuration
Open the `config.yml` file in the project root:
```yaml
target_url: "www.example.com"
number_measurements: 2
seconds_between_measurements: 2s
limit_per_measurement: 100
locations: ["DE", "Belgium", "NL"]
```
### Configuration Options

| Option | Description | Example |
|--------|-------------|---------|
| `target_url` | The website or API endpoint to test | `"www.fh-aachen.de"` |
| `number_measurements` | Number of measurement rounds to perform | `2` |
| `seconds_between_measurements` | Delay between measurement rounds | `2s` |
| `limit_per_measurement` | Maximum number of probes per measurement | `100` |
| `locations` | Array of locations to test from | See below |

### Location Formats
[Globalping](https://globalping.io) supports multiple "magic" input formats for locations:
```yaml
locations:
  - "FR"                    # Country code
  - "Poland"                # Country name
  - "Berlin+Germany"        # City + Country
  - "California"            # State/Region
  - "Europe"                # Continent
  - "Western Europe"        # Region
  - "AS13335"               # Autonomous System Number (e.g., Cloudflare)
  - "aws-us-east-1"         # Cloud provider region
  - "Google"                # Cloud provider
```
## Usage
### Running the Tool
1. Set your Globalping API token as an environment variable:
```bash
export GLOBALPING_TOKEN="Put_Your_Global_ping_token_here"
```
or in a .env file
```.env
GLOBALPING_TOKEN="Put_Your_Global_ping_token_here"
```
2. Run the tool:
```bash
make run
```
3. Review the configuration and start the measurements:
```
s: start, q: quit
(s/q): s
```
### Output
The tool generates two types of output files in the `results/` directory:
1. **Raw Data** (`.jsonl`): Complete measurement responses from the Globalping API
2. **Statistics** (`.csv`): Analyzed latency statistics for easy comparison
Example CSV output includes:
- Location information
- Minimum, maximum, average and median latencies
### Example Session
```
Config:
TargetUrl: www.fh-aachen.de
Number_measurements: 2
Seconds_between_measurements: 2s
Globalping_token: set
LimitPerMeasurement: 100
Locations: [DE Belgium NL]

Globalping Limits:
Measurements Create Limit: 500,
Remaining: 270,
Reset: 2232 seconds,
Credits Remaining: 52694

Total Limits: 52964
Required Limits/number of requests: 200

s: start, q: quit
(s/q): s

Making Measurements:
Measurements completed: 2/2 | Total duration: 9s
Number of unfinished results: 0

Detailed statistics saved to: ./results/www-fh-aachen-de_2025-12-15_18-59-08_stats.csv
```
### Unfinished Results
Some probes may not respond in time. The tool reports the count of unfinished measurements. This is normal for distributed testing.
## Output Files
Results are saved with timestamps in the `results/` directory:
```
results/
├── www-example-com_2025-12-15_18-59-08_stats.csv
└── www-example-com_2025-12-15_18-59-08.jsonl
```
## Best Practices
Run the tools for a few hours or even days to get reliable data.
## Support
For issues, questions, or feature requests, please open an issue on GitHub.

[MIT LICENSE](LICENSE)
