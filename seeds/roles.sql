/* Roles */
INSERT  into rw_roles (name,slug,created_at,updated_at) values ("Administrator","admin",current_time ,current_time );

/* Permissions */
INSERT INTO rw_permissions (name,created_at,updated_at) VALUES ("READ",current_time ,current_time );
INSERT INTO rw_permissions (name,created_at,updated_at) VALUES ("WRITE",current_time ,current_time);
INSERT INTO rw_permissions (name,created_at,updated_at) VALUES ("DELETE",current_time ,current_time);

/* Roles-Permissions */

SELECT @role_id:=id from rw_roles WHERE slug = "admin";
SELECT @read_permission_id:=id from rw_permissions where name = "READ";
SELECT @write_permission_id:=id from rw_permissions where name = "WRITE";
SELECT @delete_permission_id:=id from rw_permissions where name = "DELETE";
INSERT INTO rw_roles_permissions (created_at,updated_at,role_id,permission_id) VALUES (current_time ,current_time, @role_id,@read_permission_id);
INSERT INTO rw_roles_permissions (created_at,updated_at,role_id,permission_id) VALUES (current_time ,current_time, @role_id,@write_permission_id);
INSERT INTO rw_roles_permissions (created_at,updated_at,role_id,permission_id) VALUES (current_time ,current_time, @role_id,@delete_permission_id);