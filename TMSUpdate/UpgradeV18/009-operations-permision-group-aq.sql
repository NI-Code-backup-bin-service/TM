--multiline
CREATE TABLE `operations_permissiongroup_acquirer` (
  `permissiongroup_id` int(11) NOT NULL,
  `acquirer_profile_id` int(11) NOT NULL,
  `acquirer_name` varchar(255) DEFAULT NULL
);
