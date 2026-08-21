package com.ashhforrd.isolation.experiment;

import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequestMapping("/api/experiments")
public class IsolationExperimentController {

    private final IsolationExperimentService service;

    public IsolationExperimentController(
            IsolationExperimentService service
    ) {
        this.service = service;
    }

    @PostMapping("/dirty-read")
    public ExperimentResult runDirtyRead(
            @RequestParam IsolationLevel isolation
    ) {
        return service.runDirtyRead(isolation);
    }

    @PostMapping("/non-repeatable-read")
    public ExperimentResult runNonRepeatableRead(
            @RequestParam IsolationLevel isolation
    ) {
        return service.runNonRepeatableRead(isolation);
    }

    @PostMapping("/phantom-read")
    public ExperimentResult runPhantomRead(
            @RequestParam IsolationLevel isolation
    ) {
        return service.runPhantomRead(isolation);
    }
}