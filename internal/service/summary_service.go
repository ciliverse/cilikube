package service

import (
	"bufio"
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"k8s.io/client-go/kubernetes"

	"github.com/ciliverse/cilikube/internal/models"
	// k8s imports ... (keep existing ones)
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// Import for robust mod file parsing (optional but recommended)
	"golang.org/x/mod/modfile"
)

// SummaryService struct no longer holds client field
type SummaryService struct {
	// No longer need client kubernetes.Interface field
}

func NewSummaryService() *SummaryService {
	return &SummaryService{}
}

// GetResourceSummary Existing GetResourceSummary function ...
func (s *SummaryService) GetResourceSummary(clientSet kubernetes.Interface) (*models.ResourceSummary, map[string]error) {
	// ... (keep existing implementation) ...
	summary := &models.ResourceSummary{}
	errors := make(map[string]error)
	var wg sync.WaitGroup
	var mu sync.Mutex
	listOptions := metav1.ListOptions{Limit: 1}
	ctx := context.TODO()
	// ... (fetch funcs map and execution) ...
	fetchFuncs := map[string]func(){
		"nodes": func() {
			list, err := clientSet.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["nodes"] = err
				log.Printf("Error listing nodes: %v", err)
			} else {
				count := len(list.Items)
				summary.Nodes = &count
			}
		},
		"namespaces": func() {
			list, err := clientSet.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["namespaces"] = err
				log.Printf("Error listing namespaces: %v", err)
			} else {
				count := len(list.Items)
				summary.Namespaces = &count
			}
		},
		// ... other resource counts ...
		"pods": func() {
			list, err := clientSet.CoreV1().Pods("").List(ctx, listOptions)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["pods"] = err
				log.Printf("Error listing pods: %v", err)
			} else {
				count := len(list.Items)
				if list.RemainingItemCount != nil {
					count += int(*list.RemainingItemCount)
				} else {
					fullList, err := clientSet.CoreV1().Pods("").List(ctx, metav1.ListOptions{})
					if err == nil {
						count = len(fullList.Items)
					}
				}
				summary.Pods = &count
			}
		},
		"deployments": func() {
			list, err := clientSet.AppsV1().Deployments("").List(ctx, listOptions)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["deployments"] = err
				log.Printf("Error listing deployments: %v", err)
			} else {
				count := len(list.Items)
				if list.RemainingItemCount != nil {
					count += int(*list.RemainingItemCount)
				} else {
					fullList, err := clientSet.AppsV1().Deployments("").List(ctx, metav1.ListOptions{})
					if err == nil {
						count = len(fullList.Items)
					}
				}
				summary.Deployments = &count
			}
		},
		"services": func() {
			list, err := clientSet.CoreV1().Services("").List(ctx, listOptions)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["services"] = err
				log.Printf("Error listing services: %v", err)
			} else {
				count := len(list.Items)
				if list.RemainingItemCount != nil {
					count += int(*list.RemainingItemCount)
				} else {
					fullList, err := clientSet.CoreV1().Services("").List(ctx, metav1.ListOptions{})
					if err == nil {
						count = len(fullList.Items)
					}
				}
				summary.Services = &count
			}
		},
		"persistentVolumes": func() {
			list, err := clientSet.CoreV1().PersistentVolumes().List(ctx, listOptions)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["persistentVolumes"] = err
				log.Printf("Error listing PVs: %v", err)
			} else {
				count := len(list.Items)
				if list.RemainingItemCount != nil {
					count += int(*list.RemainingItemCount)
				} else {
					fullList, err := clientSet.CoreV1().PersistentVolumes().List(ctx, metav1.ListOptions{})
					if err == nil {
						count = len(fullList.Items)
					}
				}
				summary.PersistentVolumes = &count
			}
		},
		"pvcs": func() {
			list, err := clientSet.CoreV1().PersistentVolumeClaims("").List(ctx, listOptions)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["pvcs"] = err
				log.Printf("Error listing PVCs: %v", err)
			} else {
				count := len(list.Items)
				if list.RemainingItemCount != nil {
					count += int(*list.RemainingItemCount)
				} else {
					fullList, err := clientSet.CoreV1().PersistentVolumeClaims("").List(ctx, metav1.ListOptions{})
					if err == nil {
						count = len(fullList.Items)
					}
				}
				summary.Pvcs = &count
			}
		},
		"statefulSets": func() {
			list, err := clientSet.AppsV1().StatefulSets("").List(ctx, listOptions)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["statefulSets"] = err
				log.Printf("Error listing StatefulSets: %v", err)
			} else {
				count := len(list.Items)
				if list.RemainingItemCount != nil {
					count += int(*list.RemainingItemCount)
				} else {
					fullList, err := clientSet.AppsV1().StatefulSets("").List(ctx, metav1.ListOptions{})
					if err == nil {
						count = len(fullList.Items)
					}
				}
				summary.StatefulSets = &count
			}
		},
		"daemonSets": func() {
			list, err := clientSet.AppsV1().DaemonSets("").List(ctx, listOptions)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["daemonSets"] = err
				log.Printf("Error listing DaemonSets: %v", err)
			} else {
				count := len(list.Items)
				if list.RemainingItemCount != nil {
					count += int(*list.RemainingItemCount)
				} else {
					fullList, err := clientSet.AppsV1().DaemonSets("").List(ctx, metav1.ListOptions{})
					if err == nil {
						count = len(fullList.Items)
					}
				}
				summary.DaemonSets = &count
			}
		},
		"configMaps": func() {
			list, err := clientSet.CoreV1().ConfigMaps("").List(ctx, listOptions)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["configMaps"] = err
				log.Printf("Error listing ConfigMaps: %v", err)
			} else {
				count := len(list.Items)
				if list.RemainingItemCount != nil {
					count += int(*list.RemainingItemCount)
				} else {
					fullList, err := clientSet.CoreV1().ConfigMaps("").List(ctx, metav1.ListOptions{})
					if err == nil {
						count = len(fullList.Items)
					}
				}
				summary.ConfigMaps = &count
			}
		},
		"secrets": func() {
			list, err := clientSet.CoreV1().Secrets("").List(ctx, listOptions)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["secrets"] = err
				log.Printf("Error listing Secrets: %v", err)
			} else {
				count := len(list.Items)
				if list.RemainingItemCount != nil {
					count += int(*list.RemainingItemCount)
				} else {
					fullList, err := clientSet.CoreV1().Secrets("").List(ctx, metav1.ListOptions{})
					if err == nil {
						count = len(fullList.Items)
					}
				}
				summary.Secrets = &count
			}
		},
		"ingresses": func() {
			// Check if networking.k8s.io/v1 is available
			// For simplicity, assuming v1 is used. Add checks/fallback if needed.
			list, err := clientSet.NetworkingV1().Ingresses("").List(ctx, listOptions)
			mu.Lock()
			defer mu.Unlock()
			if err != nil {
				errors["ingresses"] = err
				log.Printf("Error listing Ingresses: %v", err)
			} else {
				count := len(list.Items)
				if list.RemainingItemCount != nil {
					count += int(*list.RemainingItemCount)
				} else {
					fullList, err := clientSet.NetworkingV1().Ingresses("").List(ctx, metav1.ListOptions{})
					if err == nil {
						count = len(fullList.Items)
					}
				}
				summary.Ingresses = &count
			}
		},
		// Add more functions for other resource types here

		// ... add other resource fetch funcs from previous example ...
	}
	wg.Add(len(fetchFuncs))
	for _, fn := range fetchFuncs {
		go func(f func()) { defer wg.Done(); f() }(fn)
	}
	wg.Wait()
	return summary, errors
}

// --- New Function to get Backend Dependencies ---

// BackendDependency represents a single Go module dependency.
type BackendDependency struct {
	Path    string `json:"path"`    // Module path (e.g., github.com/gin-gonic/gin)
	Version string `json:"version"` // Module version (e.g., v1.9.1)
}

// GetBackendDependencies reads and parses the go.mod file to list direct dependencies.
func (s *SummaryService) GetBackendDependencies() ([]BackendDependency, error) {
	// Find go.mod relative to the running executable or use a known path
	// This assumes go.mod is at the project root relative to where the binary runs.
	// Adjust this path if your structure is different or use env vars.
	modPath := "go.mod" // Adjust if needed, e.g., get from os.Getwd() or config
	if _, err := os.Stat(modPath); os.IsNotExist(err) {
		// Attempt to find it relative to executable path (might be needed in containers)
		exePath, err := os.Executable()
		if err == nil {
			modPath = filepath.Join(filepath.Dir(exePath), "go.mod")
			if _, err := os.Stat(modPath); os.IsNotExist(err) {
				log.Printf("go.mod not found at default path or executable path: %s", filepath.Join(filepath.Dir(exePath), "go.mod"))
				return nil, err // Or return empty list?
			}
		} else {
			log.Printf("go.mod not found at default path and failed to get executable path: %v", err)
			return nil, err // Or return empty list?
		}
	}

	goModBytes, err := os.ReadFile(modPath)
	if err != nil {
		log.Printf("Error reading go.mod file at %s: %v", modPath, err)
		return nil, err
	}

	modFile, err := modfile.Parse(modPath, goModBytes, nil)
	if err != nil {
		log.Printf("Error parsing go.mod file at %s: %v", modPath, err)
		// Fallback to simpler parsing if modfile fails (less robust)
		return parseGoModSimple(goModBytes)
	}

	dependencies := []BackendDependency{}
	for _, req := range modFile.Require {
		// We only care about direct dependencies (not Indirect)
		if !req.Indirect {
			dependencies = append(dependencies, BackendDependency{
				Path:    req.Mod.Path,
				Version: req.Mod.Version,
			})
		}
	}

	return dependencies, nil
}

// Simple fallback parser if golang.org/x/mod/modfile fails or is not used
func parseGoModSimple(goModBytes []byte) ([]BackendDependency, error) {
	dependencies := []BackendDependency{}
	scanner := bufio.NewScanner(strings.NewReader(string(goModBytes)))
	inRequireBlock := false
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "require (" {
			inRequireBlock = true
			continue
		}
		if line == ")" && inRequireBlock {
			inRequireBlock = false
			continue
		}
		if inRequireBlock && line != "" && !strings.HasPrefix(line, "//") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				// Check for // indirect comment
				isIndirect := false
				if len(parts) > 2 && parts[len(parts)-1] == "indirect" && parts[len(parts)-2] == "//" {
					isIndirect = true
				}
				if !isIndirect {
					dependencies = append(dependencies, BackendDependency{
						Path:    parts[0],
						Version: parts[1],
					})
				}
			}
		} else if strings.HasPrefix(line, "require") && !inRequireBlock { // Single line require
			parts := strings.Fields(line)
			if len(parts) == 3 { // require module version
				dependencies = append(dependencies, BackendDependency{
					Path:    parts[1],
					Version: parts[2],
				})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("Error scanning go.mod content: %v", err)
		return nil, err
	}
	return dependencies, nil
}
