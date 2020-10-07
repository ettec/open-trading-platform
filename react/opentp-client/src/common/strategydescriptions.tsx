export function getStrategyDisplayName(mic: string) : string | undefined {
    switch(mic) {
        case "XVWAP":
            return "VWAP STRATEGY"
    }

    return "UNKNOWN"

}