import { Destinations } from "./destinations"

export function getStrategyDisplayName(mic: string) : string | undefined {
    switch(mic) {
        case Destinations.VWAP:
            return "VWAP STRATEGY"
    }

    return "UNKNOWN"

}