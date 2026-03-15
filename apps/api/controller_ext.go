package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ashucfx/distributed-job-platform/packages/database"
)

// Extending the controller locally to support listing for the dashboard
// Since we are moving quickly, this avoids adding to the interface first.
func (c *JobController) ListJobs(ctx *gin.Context) {
	// Quick hacky implementation to serve the dashboard. 
	// In production, interface is preferred, but we have the Gorm struct.
	// Since we defined `db` as `DBService` which doesn't have List, we will cast it.
	
	if gormSvc, ok := c.db.(*database.GormDBService); ok {
		var jobs []database.Job
		// Use DB directly for list
		res := gormSvc.GetDB().Order("created_at desc").Limit(20).Find(&jobs)
		if res.Error != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list"})
			return
		}
		ctx.JSON(http.StatusOK, jobs)
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "unsupported db operation"})
	}
}
