#!/bin/bash
source /home/SCRIPTS/__lib__/common.sh
source /home/SCRIPTS/__lib__/yaml.sh
source /home/SCRIPTS/__lib__/ssh.sh
source /home/SCRIPTS/__lib__/postgresql.sh
source /home/SCRIPTS/__lib__/xo.v1.sh
source /home/SCRIPTS/__lib__/go.sh

# Get YAML file name
YAML="$1"

# Get Variables from YAML file
eval "$(parseYAML "${YAML}")"

if [[ $config_db_ssh_port ]]; then
  runSshTunnel "${config_db_ssh_port}" "${config_db_ssh_key}" "${config_db_ssh_user}" "${config_db_ssh_host}" "${config_db_port}" "${config_db_ssh_tunnel_db_host}" "${config_db_ssh_tunnel_db_port}"
fi

# Connect string to Database
CS="$(xoPostgresqlConnectString "${config_db_user}" "${config_db_pass}" "${config_db_host}" "${config_db_port}" "${config_db_name}" "${config_db_sslmode}")"
CSUID="$(xouidPostgresqlConnectString "${config_db_user}" "${config_db_pass}" "${config_db_host}" "${config_db_port}" "${config_db_name}" "${config_db_sslmode}" "${config_db_pool_max_conns}")"

runXO "${CS}" "${CSUID}" "${config_codegen_path}" "${config_codegen_ignore_fields}" "${config_codegen_package}" "${config_codegen_templates}" "${config_codegen_queries}"

replaceInterfaceToAny "${config_codegen_path}"
glueXoXouid "${config_codegen_path}"
extractRepo "${config_codegen_path}" "${config_codegen_package}"
removeXoXouid "${config_codegen_path}"
cleanXoXouidSourceBlocks "${config_codegen_path}"

goFormatCode "${config_codegen_path}"
goLintCode "${config_codegen_path}"

echo "XO finished success"
