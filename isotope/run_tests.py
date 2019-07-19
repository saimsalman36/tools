import argparse
import logging

from runner import cluster, config as cfg, consts, entrypoint, mesh, pipeline


def main(args: argparse.Namespace) -> None:
    log_level = getattr(logging, args.log_level)
    logging.basicConfig(level=log_level, format='%(levelname)s\t> %(message)s')

    config = cfg.from_toml_file(args.config_path)

    if args.clean_up == 'True':
        cluster.clean_up(
            config.cluster_project_id, config.cluster_name, config.cluster_zones)
        return

    cluster.set_up_if_not_exists(
        config.cluster_project_id, config.cluster_name, config.cluster_zones,
        config.cluster_version, config.server_machine_type,
        config.server_disk_size_gb, config.server_num_nodes,
        config.client_machine_type, config.client_disk_size_gb)

    for topology_path in config.topology_paths:
        for env_name in config.environments:
            entrypoint_service_name = entrypoint.extract_name(topology_path)
            mesh_environment = mesh.for_state(
                env_name, entrypoint_service_name,
                consts.SERVICE_GRAPH_NAMESPACE, config, args.helm_values)
            pipeline.run(topology_path, mesh_environment, config.server_image,
                         config.client_image, config.istio_archive_url,
                         config.policy_files, config.client_qps,
                         config.client_duration, config.client_num_conc_conns,
                         config.labels())


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser()
    parser.add_argument('config_path', type=str)
    parser.add_argument('helm_values', type=str)
    parser.add_argument('--clean-up',
                        type=str,
                        choices=['True', 'False'],
                        default='False')
    parser.add_argument(
        '--log_level',
        type=str,
        choices=['CRITICAL', 'ERROR', 'WARNING', 'INFO', 'DEBUG'],
        default='DEBUG')
    return parser.parse_args()


if __name__ == '__main__':
    args = parse_args()
    main(args)
