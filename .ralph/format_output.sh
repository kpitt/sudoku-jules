#!/bin/bash

# Define colors for terminal output
CYAN='\033[36m'
YELLOW='\033[33m'
GREEN='\033[32m'
RED='\033[31m'
GRAY='\033[90m'
NC='\033[0m' # No Color

MESSAGE=""
while read -r line; do

  # Check if the line is valid JSON before trying to parse it
  if echo "$line" | jq -e . >/dev/null 2>&1; then

    # Extract the event type from the JSON stream
    EVENT_TYPE=$(echo "$line" | jq -r '.type // empty')
    if [[ "$EVENT_TYPE" != "message" && -n "$MESSAGE" ]]; then
      # force a new line after last chunk of message
      echo ""
      MESSAGE=""
    fi

    case "$EVENT_TYPE" in
    "message")
      ROLE=$(echo "$line" | jq -r '.role')
      if [ "$ROLE" == "assistant" ]; then
        if [[ -z "$MESSAGE" ]]; then
          echo -e "${CYAN}[MESSAGE]${NC}"
        fi

        MSG=$(echo "$line" | jq -r '.content // empty')
        MESSAGE="${MESSAGE}${MSG}"
        echo -n "$MSG"
      fi
      ;;

    "tool_use")
      # e.g., Agent runs a shell command or reads a spec
      TOOL=$(echo "$line" | jq -r '.tool_name')
      PARAMS=$(echo "$line" | jq -r '.parameters')
      CMD=$(echo "$PARAMS" | jq -r '.command // .file_path // .dir_path // .pattern // .request // empty')
      echo -e "${YELLOW}[TOOL: $TOOL]${NC} $CMD"
      ;;

    "tool_result")
      # e.g., The result of the tests
      STATUS=$(echo "$line" | jq -r '.status // empty')
      OUTPUT=$(echo "$line" | jq -r '.output // empty')
      if [ "$STATUS" == "success" ]; then
        echo -e "${GREEN}[TOOL SUCCEEDED]${NC} $OUTPUT"
      else
        echo -e "${RED}[TOOL FAILED]${NC} $OUTPUT"
      fi
      ;;

    "result")
      STATUS=$(echo "$line" | jq -r '.status // empty')
      if [ "$STATUS" == "success" ]; then
        echo -e "${GREEN}[LOOP SUCCEEDED]${NC}"
      else
        MSG=$(echo "$line" | jq -r '.error?.message // empty')
        echo -e "${RED}[LOOP FAILED]${NC} $MSG"
      fi
      STATS=$(echo "$line" | jq -c '.stats')
      TOKENS=$(echo "$STATS" | jq -r '.total_tokens')
      echo -e "${YELLOW}[STATS] Tokens used: $TOKENS${NC}"
      echo "$STATS" | jq -r '.models | to_entries | .[] | "  Model: \(.key), \(.value.total_tokens) tokens"'
      ;;

    *)
      # Silently ignore other JSON chatter (like raw LLM text generation)
      # unless you specifically want to log it to a file.
      ;;
    esac
  fi
done
