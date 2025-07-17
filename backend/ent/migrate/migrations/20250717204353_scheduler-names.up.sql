-- Rename a column from "name" to "schedule_name"
ALTER TABLE "queue_scheduled_tasks" RENAME COLUMN "name" TO "schedule_name";
-- Modify "queue_scheduled_tasks" table
ALTER TABLE "queue_scheduled_tasks" DROP CONSTRAINT "queue_scheduled_tasks_name_key", ADD COLUMN "task_name" character varying NOT NULL;
-- Create index "queue_scheduled_tasks_schedule_name_key" to table: "queue_scheduled_tasks"
CREATE UNIQUE INDEX "queue_scheduled_tasks_schedule_name_key" ON "queue_scheduled_tasks" ("schedule_name");
