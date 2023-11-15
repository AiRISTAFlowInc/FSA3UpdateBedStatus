# A3UpdateBedStatus

This activity updates a bed to Available.

What this activity is doing:

1. Button press provides MAC
2. api/Device/GetByMacAddress for ItemId
3. api/Staff/GetByStaffId using ItemId for BedStatus
4. Then update bed Status: cd update
   - api/Staff/ChangeItemAssociation
   - api/Staff/EndItemAssociation

### Input:

| Name       | Type   | Description                                           |
| :--------- | :----- | :---------------------------------------------------- |
| IP         | string | IP address followed by port. e.g. "32.41.13.112:322"  |
| CustomerId | string | Customer unique identifier. e.g. "43"                 |
| Username   | string | Username used to log into Sofia. e.g. "adminUser"     |
| Password   | string | Password used to log into Sofia. e.g. "adminPassword" |
| MAC        | string | Bed device MAC address. e.g. "AB:CE:AF:12:00:11"      |

### Output:

| Name      | Type    | Description                      |
| :-------- | :------ | :------------------------------- |
| Status    | boolean | True if successful, false if not |
| BedStatus | string  | Holds updated bed status         |
