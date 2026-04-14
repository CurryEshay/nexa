# Nexa command guide

---

Nexa Guide (Power Users)

This is a full command reference for Nexa.

If you're new, read the README first.

Example workflow:

` mk;Education/Math/Exercise;1;2_d
  mk;Education/Physics/Worksheet;2;1_d
`

→ select Education
→ tasks are ranked automatically
→ start with top 3

## Task Architecture
Nexa uses a Project / Category / Task architecture. This means each project contains categories, which contain tasks. Tasks cannot belong directly to a project, they must belong to a category. Projects, categories and tasks are all different objects.
Paths
Because of the PCT (Project, Category, Task) architecture, Nexa uses paths to make interacting with objects easy. A path looks like this:

[Project Name]/[Category Name]/[Task Name]



In the following PCT tree, here is what some paths would look like.

Education (Education)

	Math (Education/Math)
		Task1 (Education/Math/Task1)

	English (Education/English)
		Task1 (Education/English/Task1)

	Physics (Education/Physics)
		task1 (Education/Physics/task1)

Paths can also be separated using “-”. E.g. Education-Physics-task1

## Views
Nexa displays 2 views. On the left you have a tree view of your projects and categories like before.

### Tree view example

Education
	Math
	English
	Physics
	Chemistry

1 project named Education with 4 categories: Math, English, Physics, Chemistry.

### Task view example

Task Name	            Priority        Due Date			        Days Remaining
Finish exercise 6A	P1	        Friday, Apr 15 11:59PM          2 days
Finish exercise 6B	P2	        Friday, Apr 17 11:59PM          4 days

Note: Repeating tasks will show how often they repeat

Here is an example from the app
<img width="1351" height="337" alt="5" src="https://github.com/user-attachments/assets/2984dd03-ad2a-47bc-9876-d52a299fa1ee" />

### Reference characters
Reference characters (“*” or “.”) can be used to simplify paths. A reference character represents the current location of the selected object.

"." or "*" refers to the currently selected project

"./Category/Task" refers to a task within the current project


### Complex examples

“././.” represents the current task that is selected, “././Task1”, represents Task1 in the current project and category selected. Just “.” or “*” represents the current project that is selected.

## Navigation 

Shift + Up / Down → navigate tree view  
Shift + Left / Right → navigate task view

This changes the selected object and hence changes reference operators (“*" and ".”).
Dynamic ranking
In the image above the tasks are listed in a specific order. This is the order you should do them. The tasks that show up on the right depend on the category or project selected on the left. If the selected object (highlighted in blue on the left, in this case “Nexa” in the image) is a category, all the tasks in that category appear on the right. If you have selected a project (E.g. Coding, Education or Health in picture), all tasks under that project (including all tasks in the project’s categories) appear dynamically ranked. This is how you choose what to do next. While this ranking won’t be 100% perfect (Can’t know your exact preference), a general rule of thumb is, to start with pick one in the top 3.

Goal: reduce decision-making.
You should not scan the full list. Start with the top 3.

## Locking
To avoid mistakes, to make any sort of deletion or moving, the program must be unlocked. By default the program is locked, to unlock type “ulk”, to lock again type “lk”. You can still complete tasks while the project is locked.

## Commands
A command has a key word followed by arguments. A key word calls a function, meaning a keyword invokes a functionality, the arguments are the data required to execute a function. For example a remove function would need the path of the object you want to remove. The keywords and arguments are separated by semicolons.

## Make
Makes an object.

`
mk;[path/to/object];[extra args when making a task]
`
### Make project

`
mk;[path/to/project]
`

E.g. mk;myProject

### Make categories

`
mk;[path/to/category]
`

E.g. mk;Education/Math

Note: For this to work the “Education” project must already exist.

### Make task

`
mk;[path/to/task];[priority];[date];[time - optional, default 11:59PM]
`

E.g. mk;Education/Math/Exercise;1;3_d;15:00

Note: For this to work, the Education project and Math category must exist. Note a short date is here. The syntax is shown below. An absolute date in DD-MM-YYYY can also be used. 12 hour time instead of 24 hour time can also be used.

E.g. mk;Education/Math/Exercise;1;10/01/2027;3:00PM

### Short date syntax
x_d: x days from today
x_w: x weeks from today
x_m: x months from today
x_y: x years from today

## Remove
Removes an object. This is recursive, meaning if you delete a category it will delete all the tasks inside, or if you delete a project it will delete all the categories and tasks inside it. THIS CANNOT BE UNDONE. To prevent mistakes a locking feature has been added, if the program is locked, no removal can be done. To unlock type “ulk”, to lock again type “lk”.

`
rm;[path/to/object]
`

### Remove project

`
rm;[path/to/project]
`

E.g. rm;myProject

### Remove categories

`
rm;[path/to/category]
`

E.g. rm;Education/Math
Note: For this to work the “Education” project must already exist.

### Remove task

`
rm;[path/to/task]
`

E.g. rm;Education/Math/Exercise

Note: For this to work, the Education project and Math category must exist.

## Move
Moves an object to a different path, can be used for renaming. The program must be unlocked to use this function.

`
mv;[old_path/to/old_object];[new_path/to/new_object]
`

### Move project (Essentially renaming)

`
mv;[old_path/to/project];[new_path/to/project]
`

E.g. mv;myProject;myProjectNew

### Move categories

`
mv;[old_path/to/category];[new_path/to/category]
`

Renaming
E.g. mv;Education/Math;Education/English

Moving
E.g. mv;Education/Math;newEducation/English

Note: For this to work the “Education” and “newEducation” project must already exist.

### Move task

`
mv;[old_path/to/task];[new_path/to/task]
`

Moving
E.g. mv;Education/Math/Exercise;Education/English/Exercise2

Renaming
E.g. mv;Education/Math/Exercise;Education/Math/Exercise2

Note: For this to work, the Education project and Math and English categories must exist.

## Repeating tasks
Nexa supports repeating tasks, for example workout every friday. The syntax changes if you want to create a repeating task.

`
rp;[path/to/task];[priority];[increment duration];[reference date];[time - optional, default 11:59PM]
`

The first 2 arguments are the same as making a one time task, although there are some differences.

### Increment duration

[increment duration] - The space between each occurrence of the task. E.g. workout every Tuesday, the space between each task is 1 week. The way you input this looks like this.

{years} {months} {weeks} {days}

Note: These are delimited by a space. Here are some examples.

Every week: 0 0 1 0		Every month: 0 1 0 0		Every month and 1 day: 0 0 1 1

Every 8 days: 0 0 1 1		Every year: 1 0 0 0		Every 2 years and 5 months: 2 5 0 0

### Reference date
[reference date] - Acts as a reference so Nexa knows when to repeat. For working out every Tuesday, this would be the date reference to the last Tuesday or the coming Tuesday. Basically, what is one example of this task occurring? 

Note: You may use short dates for this.


E.g. Running every Wednesday at 4pm, suppose today is a Monday (13/04/2026)

rp;Health/Running/Run 2km;1;0 0 1 0;2_d;4:00pm

Or 

rp;Health/Running/Run 2km;1;0 0 1 0;13/04/2026;4:00pm


## Done task
Used to complete a task. Does not require unlocking of the program. Do this to complete a task instead of remove (rm). Works for both one off and repeating tasks.

`
dn;[path/to/task] 
`

E.g.

dn;Education/CS/Finish assignment



## Reschedule
Used to reschedule a task. Works for repeating and one-off tasks.

`
rs;[path/to/task];[date];[time - optional, default 11:59PM]
`

E.g.

rs;Education/CS/Assignment;1_m;13:00



## Reprioritise 
Used to change the priority of a task. Works for repeating and one-off tasks.

`
p;[path/to/task];[priority]
`

E.g.

p;Education/CS/Assignment;3
