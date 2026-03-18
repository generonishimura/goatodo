export namespace presenter {
	
	export class TaskResponse {
	    success: boolean;
	    data?: any;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new TaskResponse(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.data = source["data"];
	        this.error = source["error"];
	    }
	}
	export class UpdateTaskRequest {
	    id: string;
	    title?: string;
	    status?: string;
	    priority?: number;
	
	    static createFrom(source: any = {}) {
	        return new UpdateTaskRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.status = source["status"];
	        this.priority = source["priority"];
	    }
	}

}

