package cmd

import (
	"cmd-ai-resolver/internal/llm"
	"cmd-ai-resolver/internal/logger"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type commandHandler struct {
	filePath        string
	fileContent     string
	modifiedContent string
	aiTagRegexp     *regexp.Regexp
}

func newCommandHandler(filePath string) *commandHandler {
	return &commandHandler{
		filePath:        filePath,
		fileContent:     "",
		modifiedContent: "",
		aiTagRegexp:     regexp.MustCompile(`<AI>(.*?)</AI>`),
	}
}

func (cm *commandHandler) handleCommand() error {
	logger.Log.Debugf("Processing file: %s", cm.filePath)

	err := cm.readFileContent()
	if err != nil {
		logger.Log.Errorf("Error reading file %s: %v", cm.filePath, err)
		return fmt.Errorf("reading file: %w", err)
	}
	logger.Log.Debugf("Original content:\n%s", cm.fileContent) // Using Debugf for visibility with -v

	originalTag, extractedPrompt := cm.extractAITag()

	if extractedPrompt != "" {
		logger.Log.Debugf("Extracted AI prompt: %s", extractedPrompt)
	} else {
		return cm.passThroughCommandExec()
	}

	processedSegment := ""
	if extractedPrompt != "" {
		logger.Log.Debugf("Preparing to process with LLM...")

		segment, llmErr := llm.ProcessWithOpenAI(cm.fileContent, extractedPrompt)
		if llmErr != nil {
			logger.Log.Errorf("Error processing with LLM: %v", llmErr)
			return fmt.Errorf("processing with LLM: %w", llmErr)
		}
		processedSegment = segment
		logger.Log.Debugf("LLM processing finished. Received segment: %s", processedSegment)
	}

	cm.modifiedContent = cm.fileContent
	madeChanges := false
	if originalTag != "" && processedSegment != "" {
		cm.modifiedContent = strings.Replace(cm.fileContent, originalTag, processedSegment, 1)
		logger.Log.Debugf("Modified content:\n%s", cm.modifiedContent)
		madeChanges = true
	} else if originalTag != "" && extractedPrompt != "" {
		logger.Log.Debugf("AI tag was found but no processed segment was generated. Original content will be preserved for the tag.")
	} else {
		logger.Log.Debugf("Content remains unchanged.")
	}

	if madeChanges {
		err = cm.writeProcessedFile()
		if err != nil {
			logger.Log.Errorf("Error writing modified file %s: %v", cm.filePath, err)
			return fmt.Errorf("writing modified file: %w", err)
		}
	} else {
		logger.Log.Debugf("No changes made to the file.")
	}

	return nil
}

func (cm *commandHandler) readFileContent() error {
	content, err := os.ReadFile(cm.filePath)
	if err != nil {
		return err
	}
	cm.fileContent = string(content)
	return nil
}

func (cm *commandHandler) extractAITag() (originalTag string, extractedPrompt string) {
	matches := cm.aiTagRegexp.FindStringSubmatch(cm.fileContent)

	if len(matches) == 2 {
		originalTag = matches[0]
		extractedPrompt = matches[1]
	}
	return originalTag, extractedPrompt
}

func (cm *commandHandler) writeProcessedFile() error {
	err := os.WriteFile(cm.filePath, []byte(cm.modifiedContent), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (cm *commandHandler) passThroughCommandExec() error {
	logger.Log.Debugf("No AI tags found. Running pass-through command: %s %s", passThroughCommand, cm.filePath)

	if passThroughCommand == "" {
		logger.Log.Debugf("No pass-through command provided. Skipping execution.")
		return nil
	}
	cmd := exec.Command(passThroughCommand, cm.filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	return cmd.Run()
}
