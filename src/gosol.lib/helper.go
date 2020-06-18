package main

func destinationTypeToString(i int) string {
	switch i {
		case 1: return "TOPIC"
		case 2: return "QUEUE"
	}
	return "NONE"
}

func connectionTypeToString(i int) string {
	switch i {
		case 1: return "UP"
		case 2: return "RECONNECTING"
		case 3: return "RECONNECTED"
	}
	return "DOWN"
}

func publishedEventTypeToString(i int) string {
	switch i {
	case 1:
		return "ACK"
	}
	return "REJECT"
}
